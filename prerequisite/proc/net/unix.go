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
	return p.CheckTemplate(func() {
		content, err := os.ReadFile("/proc/net/unix")
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("reading /proc/net/unix: %w", err))
			return
		}
		p.Satisfied = strings.Contains(string(content), p.Expected)
		return
	})
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
