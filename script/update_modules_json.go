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
	Doc         string   `json:"doc,omitempty"`
	Children    []module `json:"children,omitempty"`
}

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
		SchemaVersion: "1.0",
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
	m := module{
		Name:        name,
		Aliases:     aliases,
		Description: description,
		Doc:         resolveDoc(root, kind, name),
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

func existingDoc(root, doc string) string {
	docPath := strings.SplitN(doc, "#", 2)[0]
	if docPath == "" {
		return ""
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(docPath, "./")))); err != nil {
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
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(docPath, "./")))); err != nil {
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
