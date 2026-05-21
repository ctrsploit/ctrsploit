package pipeprimitive

import (
	"bytes"
	"debug/elf"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ctrsploit/ctrsploit/pkg/pipe-primitive/runcwatch"
)

func TestShellPayload(t *testing.T) {
	got := string(ShellPayload("touch /escaped"))
	want := "#!/bin/bash\ntouch /escaped\n"
	if got != want {
		t.Fatalf("payload = %q, want %q", got, want)
	}
}

func TestStaticShellPayload(t *testing.T) {
	payload, err := StaticShellPayload("touch /escaped")
	if err != nil {
		t.Fatalf("StaticShellPayload returned error: %v", err)
	}
	if !bytes.HasPrefix(payload, []byte("\x7fELF")) {
		t.Fatalf("static shell payload does not start with ELF magic: % x", payload[:4])
	}
	for _, want := range []string{"/bin/bash", "-c", "touch /escaped"} {
		if !bytes.Contains(payload, []byte(want)) {
			t.Fatalf("static shell payload missing argv %q", want)
		}
	}
}

func TestStaticShellPayloadRejectsLongArg(t *testing.T) {
	_, err := StaticShellPayload(strings.Repeat("A", staticShellArgBufSize))
	if err == nil {
		t.Fatal("expected oversized arg error")
	}
}

func TestStaticShellPayloadAMD64IsSmallELF(t *testing.T) {
	payload, err := StaticShellPayload("true")
	if err != nil {
		t.Fatalf("StaticShellPayload returned error: %v", err)
	}
	if len(payload) > 4096 {
		t.Fatalf("static shell payload length = %d, want <= 4096", len(payload))
	}

	f, err := elf.NewFile(bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("parse static shell payload ELF: %v", err)
	}
	defer f.Close()
	if f.Type != elf.ET_EXEC {
		t.Fatalf("static shell payload type = %s, want ET_EXEC", f.Type)
	}
	if f.Machine != elf.EM_X86_64 {
		t.Fatalf("static shell payload machine = %s, want EM_X86_64", f.Machine)
	}
}

func TestStaticShellPayloadExecutesCommand(t *testing.T) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		t.Skipf("static shell payload execution unsupported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	dir := t.TempDir()
	proof := filepath.Join(dir, "proof")
	payload, err := StaticShellPayload("touch " + proof)
	if err != nil {
		t.Fatalf("StaticShellPayload returned error: %v", err)
	}
	path := filepath.Join(dir, "payload")
	if err := os.WriteFile(path, payload, 0o755); err != nil {
		t.Fatalf("write payload: %v", err)
	}
	if out, err := exec.Command(path).CombinedOutput(); err != nil {
		t.Fatalf("execute static shell payload: %v\n%s", err, out)
	}
	if _, err := os.Stat(proof); err != nil {
		t.Fatalf("payload did not create proof file: %v", err)
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
