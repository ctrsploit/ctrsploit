package pipeprimitive

import (
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
