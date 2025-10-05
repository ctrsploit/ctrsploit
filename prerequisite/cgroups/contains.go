package cgroups

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type Contains struct {
	prerequisite.BasePrerequisite
	Expected string
}

func (p *Contains) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	file, err := os.Open("/proc/1/cgroup")
	if err != nil {
		return false, fmt.Errorf("failed to check %s caused by reading /proc/1/cgroup: %w", p.Name, err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] != '0' {
			p.Satisfied = strings.Contains(line, p.Expected)
			if p.Satisfied {
				break
			}
		}
	}
	p.Checked = true
	return p.Satisfied, nil
}

var (
	// ContainsDocker check cgroupsPath contains /docker
	//parent default do docker when use cgroupfs
	//https://github.com/moby/moby/blob/v28.4.0/daemon/oci_linux.go#L779
	//2:cpu,cpuacct:/docker/2403ccfe427f1dc7b6120dccb9ac58c2f8456950380d587a9f33ea7f00dd303d
	//1:name=systemd:/docker/2403ccfe427f1dc7b6120dccb9ac58c2f8456950380d587a9f33ea7f00dd303d
	//0::/system.slice/containerd.service
	ContainsDocker = Contains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "/proc/1/cgroup",
			Info:   "/proc/1/cgroup contains /docker",
			ExeEnv: exeenv.InContainer,
		},
		Expected: "/docker",
	}
	// ContainsKubepods
	// https://github.com/kubernetes/kubernetes/blob/v1.34.1/pkg/kubelet/cm/node_container_manager_linux.go#L41
	// https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/cm/container_manager_linux.go#L273
	ContainsKubepods = Contains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "/proc/1/cgroup",
			Info:   "/proc/1/cgroup contains /",
			ExeEnv: exeenv.InContainer,
		},
		Expected: "kubepods",
	}
)
