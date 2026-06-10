package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	checksecCmd "github.com/ctrsploit/ctrsploit/cmd/ctrsploit/checksec"
	envCmd "github.com/ctrsploit/ctrsploit/cmd/ctrsploit/env"
	exploitCmd "github.com/ctrsploit/ctrsploit/cmd/ctrsploit/exploit"
	vulCmd "github.com/ctrsploit/ctrsploit/cmd/ctrsploit/vul"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

type catalog struct {
	SchemaVersion string  `json:"schema_version"`
	Modules       modules `json:"modules"`
}

type modules struct {
	Env      []module `json:"env"`
	Checksec []module `json:"checksec"`
	Vul      []module `json:"vul"`
	Exploit  []module `json:"exploit"`
}

type module struct {
	Name        string   `json:"name"`
	Aliases     []string `json:"aliases"`
	Description string   `json:"description"`
	Maintainer  []string `json:"maintainer"`
	Doc         string   `json:"doc,omitempty"`
	Children    []module `json:"children,omitempty"`
}

const defaultMaintainer = "ssst0n3"

var docOverrides = map[string]string{
	"env:storage-driver":      "./env/storagedriver/README.md",
	"env:no-new-privs":        "./env/nonewprivs/README.md",
	"fork-bomb":               "./vul/fork-bomb/README.md",
	"shocker":                 "./vul/caps/shocker/README.md",
	"release_agent":           "./vul/caps/sys_admin/release_agent/README.md",
	"ebpf-bash":               "./vul/caps/sys_admin/ebpf/bash/README.md",
	"ebpf-execve":             "./vul/caps/sys_admin/ebpf/execve/README.md",
	"ebpf-cron":               "./vul/caps/sys_admin/ebpf/cron/README.md",
	"ptrace-pid-host":         "./vul/caps/sys_ptrace/pid_host/README.md",
	"host-pid-proc-root":      "./vul/namespace/pid/proc_root/README.md",
	"kubeconfig-user-exec":    "./vul/kubeconfig/user-exec/README.md",
	"sa-token-access-secrets": "./vul/sa-token/access-secrets/README.md",
	"sa-token-policy":         "./vul/sa-token/policy/README.md",
	"docker.sock":             "./vul/shared-socket/docker-sock/README.md",
	"CVE-2021-22555":          "./exploit/CVE-2021-22555_ubuntu18.04/README.md",
	"crash":                   "./pkg/crash/README.md",
}

var descriptionOverrides = map[string]string{
	"fork-bomb": "fork bomb causes denial of service when resource limits or cgroup configs are unsafe",
}

func main() {
	root, err := findRepoRoot()
	if err != nil {
		fatal(err)
	}

	c := catalog{
		SchemaVersion: "1.1",
		Modules: modules{
			Env:      modulesFromCommands(root, "env", envCmd.Command.Commands, false),
			Checksec: modulesFromCommands(root, "checksec", checksecCmd.Command.Commands, false),
			Vul:      modulesFromCommands(root, "vul", vulCmd.Command.Commands, true),
			Exploit:  modulesFromCommands(root, "exploit", exploitCmd.Command.Commands, false),
		},
	}

	if err := validateDocs(root, c); err != nil {
		fatal(err)
	}
	if err := validateMaintainers(c); err != nil {
		fatal(err)
	}

	content, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fatal(err)
	}
	content = append(content, '\n')

	output := filepath.Join(root, "modules.json")
	if err := os.WriteFile(output, content, 0o644); err != nil {
		fatal(err)
	}
	fmt.Printf("Updated %s\n", output)
}

func modulesFromCommands(root, kind string, commands []*cli.Command, includeChildren bool) []module {
	result := make([]module, 0, len(commands))
	for _, cmd := range commands {
		if cmd == nil || cmd.Hidden {
			continue
		}
		result = append(result, moduleFromCommand(root, kind, cmd, includeChildren))
	}
	return result
}

func moduleFromCommand(root, kind string, cmd *cli.Command, includeChildren bool) module {
	name, aliases, description := commandMetadata(kind, cmd)
	doc := resolveDoc(root, kind, name)
	m := module{
		Name:        name,
		Aliases:     aliases,
		Description: description,
		Maintainer:  resolveMaintainers(root, doc),
		Doc:         doc,
	}
	if includeChildren {
		for _, child := range cmd.Commands {
			if child == nil || child.Hidden || isOperationCommand(child) {
				continue
			}
			m.Children = append(m.Children, moduleFromCommand(root, kind, child, false))
		}
	}
	return m
}

func commandMetadata(kind string, cmd *cli.Command) (string, []string, string) {
	if kind == "exploit" && cmd.Name == "exploit" && len(cmd.Aliases) == 1 && cmd.Aliases[0] == "x" {
		return "cve-2026-43500",
			[]string{"43500", "dirty-frag-rxrpc", "dirtyfrag-rxrpc"},
			"local privilege escalation in Linux kernel RxRPC/rxkad Dirty Frag path"
	}

	aliases := append([]string{}, cmd.Aliases...)
	if aliases == nil {
		aliases = []string{}
	}
	description := cmd.Usage
	if override, ok := descriptionOverrides[cmd.Name]; ok {
		description = override
	}
	return cmd.Name, aliases, description
}

func isOperationCommand(cmd *cli.Command) bool {
	switch cmd.Name {
	case "checksec", "exploit":
		return true
	default:
		return false
	}
}

func resolveDoc(root, kind, name string) string {
	if doc, ok := docOverrides[kind+":"+name]; ok {
		return existingDoc(root, doc)
	}
	if doc, ok := docOverrides[name]; ok {
		return existingDoc(root, doc)
	}
	lowerName := strings.ToLower(name)
	if strings.HasPrefix(lowerName, "cve-") {
		return existingDoc(root, fmt.Sprintf("./vul/%s/README.md", lowerName))
	}
	if kind == "env" {
		return existingDoc(root, fmt.Sprintf("./env/%s/README.md", name))
	}
	return ""
}

type readmeFrontMatter struct {
	Maintainer maintainerList `yaml:"maintainer"`
}

type maintainerList []string

func (m *maintainerList) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		*m = maintainerList{value.Value}
	case yaml.SequenceNode:
		var maintainers []string
		for _, item := range value.Content {
			if item.Kind != yaml.ScalarNode {
				continue
			}
			maintainers = append(maintainers, item.Value)
		}
		*m = maintainerList(maintainers)
	}
	return nil
}

func resolveMaintainers(root, doc string) []string {
	if doc == "" {
		return []string{defaultMaintainer}
	}

	maintainers, err := maintainersFromDoc(root, doc)
	if err != nil || len(maintainers) == 0 {
		return []string{defaultMaintainer}
	}
	return maintainers
}

func maintainersFromDoc(root, doc string) ([]string, error) {
	content, err := os.ReadFile(docFilepath(root, doc))
	if err != nil {
		return nil, err
	}
	return maintainersFromReadme(content)
}

func maintainersFromReadme(content []byte) ([]string, error) {
	frontMatter, ok := readFrontMatter(content)
	if !ok {
		return nil, nil
	}

	var meta readmeFrontMatter
	if err := yaml.Unmarshal(frontMatter, &meta); err != nil {
		return nil, err
	}
	return normalizeMaintainers([]string(meta.Maintainer)), nil
}

func readFrontMatter(content []byte) ([]byte, bool) {
	text := string(content)
	if !strings.HasPrefix(text, "---\n") && !strings.HasPrefix(text, "---\r\n") {
		return nil, false
	}

	lines := strings.SplitAfter(text, "\n")
	start := len(lines[0])
	offset := start
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "---" {
			return []byte(text[start:offset]), true
		}
		offset += len(line)
	}
	return nil, false
}

func normalizeMaintainers(maintainers []string) []string {
	result := make([]string, 0, len(maintainers))
	for _, maintainer := range maintainers {
		maintainer = strings.TrimSpace(maintainer)
		if maintainer == "" {
			continue
		}
		result = append(result, maintainer)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func docFilepath(root, doc string) string {
	docPath := strings.SplitN(doc, "#", 2)[0]
	return filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(docPath, "./")))
}

func existingDoc(root, doc string) string {
	docPath := strings.SplitN(doc, "#", 2)[0]
	if docPath == "" {
		return ""
	}
	if _, err := os.Stat(docFilepath(root, doc)); err != nil {
		return ""
	}
	return doc
}

func validateDocs(root string, c catalog) error {
	var missing []string
	for _, mod := range allModules(c.Modules) {
		docPath := strings.SplitN(mod.Doc, "#", 2)[0]
		if docPath == "" {
			continue
		}
		if _, err := os.Stat(docFilepath(root, mod.Doc)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				missing = append(missing, mod.Doc)
				continue
			}
			return err
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing doc paths: %s", strings.Join(missing, ", "))
	}
	return nil
}

func validateMaintainers(c catalog) error {
	var missing []string
	for _, mod := range allModules(c.Modules) {
		if len(mod.Maintainer) == 0 {
			missing = append(missing, mod.Name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing maintainers: %s", strings.Join(missing, ", "))
	}
	return nil
}

func allModules(ms modules) []module {
	var result []module
	for _, list := range [][]module{ms.Env, ms.Checksec, ms.Vul, ms.Exploit} {
		for _, mod := range list {
			result = append(result, mod)
			result = append(result, mod.Children...)
		}
	}
	return result
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("could not find repository root")
		}
		dir = parent
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
