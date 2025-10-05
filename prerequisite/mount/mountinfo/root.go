package mountinfo

import (
	"fmt"
	"strings"

	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type RootContains struct {
	prerequisite.BasePrerequisite
	MountPoint string
	Expected   string
}

func (p *RootContains) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	info, err := mountinfo.GetMountByMountpoint(p.MountPoint)
	if err != nil {
		return false, fmt.Errorf("failed to check %s caused by getting mountinfo of /etc/hosts: %w", p.Name, err)
	}
	p.Satisfied = strings.Contains(info.Root, p.Expected)
	p.Checked = true
	return p.Satisfied, nil
}

var (
	// HostsRootContainsDocker
	//https://github.com/moby/moby/blob/v28.4.0/daemon/container_operations_unix.go#L552
	//https://github.com/moby/moby/blob/v28.4.0/daemon/config/config_linux.go#L189
	//https://github.com/moby/moby/blob/v28.4.0/container/container.go#L407
	HostsRootContainsDocker = RootContains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "/etc/hosts",
			Info:   "/etc/hosts's mountinfo root contains 'docker', e.g., 814 696 259:2 /var/lib/docker/containers/44bec6602ccfae7458bcd71279beafc287d9fde509fa861377e162270d4cd92f/hosts /etc/hosts rw,relatime - ext4 /dev/nvme0n1p2 rw,errors=remount-ro",
			ExeEnv: exeenv.InContainer,
		},
		MountPoint: "/etc/hosts",
		Expected:   "docker",
	}

	// HostnameRootContainsContainerd
	// https://github.com/containerd/containerd/blob/v2.1.4/internal/cri/server/container_create.go#L1098
	// https://github.com/containerd/containerd/blob/v2.1.4/internal/cri/server/helpers.go#L86
	// https://github.com/containerd/containerd/blob/v2.1.4/defaults/defaults_unix.go#L26
	HostnameRootContainsContainerd = RootContains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "/etc/hostname",
			Info:   "/etc/hostname's mountinfo root contains 'containerd'",
			ExeEnv: exeenv.InContainer,
		},
		MountPoint: "/etc/hostname",
		Expected:   "containerd",
	}
	// HostnameRootContainsNerdctl
	// https://github.com/containerd/nerdctl/blob/v2.1.6/docs/dir.md?plain=1#L28
	HostnameRootContainsNerdctl = RootContains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "/etc/hostname",
			Info:   "/etc/hostname's mountinfo root contains 'nerdctl'",
			ExeEnv: exeenv.InContainer,
		},
		MountPoint: "/etc/hostname",
		Expected:   "nerdctl",
	}
	HostsRootContainsPods = RootContains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "/etc/hosts",
			Info:   "/etc/hosts's mountinfo root contains 'pods'",
			ExeEnv: exeenv.InContainer,
		},
		MountPoint: "/etc/hosts",
		Expected:   "pods",
	}
)
