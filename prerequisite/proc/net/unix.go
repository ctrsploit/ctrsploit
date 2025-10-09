package net

import (
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type UnixContains struct {
	prerequisite.BasePrerequisite
	Expected string
}

func (p *UnixContains) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	content, err := os.ReadFile("/proc/net/unix")
	if err != nil {
		return false, fmt.Errorf("failed to check %s, caused by reading /proc/net/unix: %w", p.Name, err)
	}
	p.Satisfied = strings.Contains(string(content), p.Expected)
	p.Checked = true
	return p.Satisfied, nil
}

var (
	// UnixContainsDockerSock
	// https://github.com/moby/moby/blob/v28.4.0/opts/hosts.go#L23
	UnixContainsDockerSock = UnixContains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "/proc/net/unix",
			Info:   "/proc/net/unix contains docker.sock, which can be seen when use host net ns",
			ExeEnv: exeenv.InContainer,
		},
		Expected: "docker.sock",
	}
	// UnixContainsContainerdSock
	// https://github.com/containerd/containerd/blob/v2.1.4/defaults/defaults_linux.go#L21
	UnixContainsContainerdSock = UnixContains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "/proc/net/unix",
			Info:   "/proc/net/unix contains containerd.sock, which can be seen when use host net ns",
			ExeEnv: exeenv.InContainer,
		},
		Expected: "containerd.sock",
	}
)
