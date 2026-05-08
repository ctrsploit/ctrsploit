package pipeprimitive

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ctrsploit/ctrsploit/pkg/pipe-primitive/runcwatch"
)

func TestCommandExploitSubcommandsUseShortNames(t *testing.T) {
	cmd := Command(&recordingPrimitive{}, nil, "recording primitive")

	topLevel := map[string]bool{}
	for _, sub := range cmd.Commands {
		topLevel[sub.Name] = true
	}
	for _, want := range []string{"privilege-escalate", "escape", "image-pollution"} {
		if !topLevel[want] {
			t.Fatalf("top-level exploit subcommand %q missing from %+v", want, topLevel)
		}
	}
	for _, unwanted := range []string{"recording-privilege-escalate", "recording-escape", "recording-image-pollution"} {
		if topLevel[unwanted] {
			t.Fatalf("top-level exploit subcommand %q should use a short name", unwanted)
		}
	}

	var escapeCmdFound bool
	for _, sub := range cmd.Commands {
		if sub.Name != "escape" {
			continue
		}
		escapeCmdFound = true

		got := map[string]bool{}
		for _, escapeSub := range sub.Commands {
			got[escapeSub.Name] = true
		}
		for _, want := range []string{"image", "exec", "restart"} {
			if !got[want] {
				t.Fatalf("escape subcommand %q missing from %+v", want, got)
			}
		}
	}
	if !escapeCmdFound {
		t.Fatal("escape command missing")
	}
}

func TestShellPayload(t *testing.T) {
	got := string(ShellPayload("touch /escaped"))
	want := "#!/bin/bash\ntouch /escaped\n"
	if got != want {
		t.Fatalf("payload = %q, want %q", got, want)
	}
}

func TestGenerateEscapeImage(t *testing.T) {
	dir := t.TempDir()
	writer := []byte("package main\nfunc main() {}\n")
	extraFiles := map[string][]byte{
		"go.mod":              []byte("module test\n"),
		"writerdep/helper.go": []byte("package writerdep\n"),
	}
	if err := GenerateEscapeImage(dir, writer, ShellPayload("touch /escaped"), extraFiles); err != nil {
		t.Fatalf("GenerateEscapeImage returned error: %v", err)
	}

	for _, name := range []string{"Dockerfile", "ld.go", "writer.go", "payload", "go.mod", filepath.Join("writerdep", "helper.go")} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("expected %s to be generated: %v", name, err)
		}
	}

	payload, err := os.ReadFile(filepath.Join(dir, "payload"))
	if err != nil {
		t.Fatalf("read payload: %v", err)
	}
	if string(payload) != "#!/bin/bash\ntouch /escaped\n" {
		t.Fatalf("payload = %q", payload)
	}
}

func TestGenerateEscapeImageExecMode(t *testing.T) {
	dir := t.TempDir()
	writer := []byte("package main\nfunc main() {}\n")
	if err := GenerateEscapeImageWithMode(dir, EscapeImageModeExec, writer, ShellPayload("touch /escaped")); err != nil {
		t.Fatalf("GenerateEscapeImageWithMode returned error: %v", err)
	}

	for _, name := range []string{"Dockerfile", "runc-capture.go", "runcwatch/runcwatch.go", "writer.go", "payload"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("expected %s to be generated: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "ld.go")); !os.IsNotExist(err) {
		t.Fatalf("ld.go should not be generated in exec mode: %v", err)
	}

	dockerfile, err := os.ReadFile(filepath.Join(dir, "Dockerfile"))
	if err != nil {
		t.Fatalf("read Dockerfile: %v", err)
	}
	content := string(dockerfile)
	for _, want := range []string{"HEALTHCHECK", "CMD [\"/proc/self/exe\"]", "ENTRYPOINT [\"/runc-capture\"]", "runc exec"} {
		if !strings.Contains(content, want) {
			t.Fatalf("Dockerfile does not contain %q:\n%s", want, content)
		}
	}
}

func TestGenerateEscapeImageRejectsUnsupportedMode(t *testing.T) {
	err := GenerateEscapeImageWithMode(t.TempDir(), "unknown", []byte("package main\nfunc main() {}\n"), ShellPayload("true"))
	if err == nil {
		t.Fatal("expected unsupported mode error")
	}
}

func TestIsRunCCmdline(t *testing.T) {
	tests := []struct {
		name    string
		cmdline string
		want    bool
	}{
		{name: "runc command", cmdline: "runc\x00init\x00", want: true},
		{name: "docker runc command", cmdline: "docker-runc\x00init\x00", want: true},
		{name: "proc self exe init", cmdline: "/proc/self/exe\x00init\x00", want: true},
		{name: "bare proc self exe", cmdline: "/proc/self/exe\x00", want: false},
		{name: "runc capture helper", cmdline: "/runc-capture\x00", want: false},
		{name: "containerd shim runc", cmdline: "containerd-shim-runc-v2\x00", want: false},
		{name: "unrelated command", cmdline: "bash\x00", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRunCCmdline(tt.cmdline); got != tt.want {
				t.Fatalf("isRunCCmdline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRunCProcess(t *testing.T) {
	tests := []struct {
		name    string
		cmdline string
		exe     string
		want    bool
	}{
		{name: "proc self exe points to runc", cmdline: "/proc/self/exe\x00", exe: "/runc", want: true},
		{name: "proc self exe points to docker runc", cmdline: "/proc/self/exe\x00", exe: "/usr/bin/docker-runc", want: true},
		{name: "proc self exe points elsewhere", cmdline: "/proc/self/exe\x00", exe: "/bin/true", want: false},
		{name: "runc cmdline wins", cmdline: "runc\x00init\x00", exe: "/bin/true", want: true},
		{name: "helper rejected", cmdline: "/runc-capture\x00", exe: "/runc-capture", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runcwatch.IsRunCProcess(tt.cmdline, tt.exe); got != tt.want {
				t.Fatalf("IsRunCProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}
