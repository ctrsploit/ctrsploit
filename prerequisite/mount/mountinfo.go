package mount

import (
	"fmt"
	"strings"

	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type RootMountInfoSourceContains struct {
	prerequisite.BasePrerequisite
	Expected string
}

func (p *RootMountInfoSourceContains) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	info, err := mountinfo.RootMount()
	if err != nil {
		return false, fmt.Errorf("failed to check rootfs's mountinfo source caused by getting root's mountinfo: %w", err)
	}
	p.Satisfied = strings.Contains(info.Source, p.Expected)
	p.Checked = true
	return p.Satisfied, nil
}

var (
	// RootMountInfoSourceContainsDocker
	// 367 342 253:1 /rootfs / rw,relatime master:152 - ext4 /dev/mapper/docker-8:1-132673-1a22af7e6fb302ee262482507d2c0247f238693a233e31a327e3c813b3c67f42 rw,stripe=16
	RootMountInfoSourceContainsDocker = RootMountInfoSourceContains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "rootfs source",
			Info:   "rootfs's mountinfo source contains 'docker', e.g., /dev/mapper/docker-8:...",
			ExeEnv: exeenv.InContainer,
		},
		Expected: "docker",
	}
)

type RootMountInfoVFSOptionsContains struct {
	prerequisite.BasePrerequisite
	Expected string
}

func (p *RootMountInfoVFSOptionsContains) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	info, err := mountinfo.RootMount()
	if err != nil {
		return false, fmt.Errorf("failed to check rootfs's mountinfo vfs options caused by getting root's mountinfo: %w", err)
	}
	p.Satisfied = strings.Contains(info.VFSOptions, p.Expected)
	p.Checked = true
	return p.Satisfied, nil
}

var (
	// RootMountInfoVFSOptionsContainsDocker
	// 1991 1838 0:108 / / rw,relatime - overlay overlay rw,lowerdir=/var/lib/docker/overlay2/l/PI4U3C2HAV5G3PKDAPQRGXYVYD:/var/lib/docker/overlay2/l/Z2EI5ICE4E4LHFITT5P5C4HEOC,upperdir=/var/lib/docker/overlay2/69df9ec3d077f103f8a9e1eab68cc60c18dcae03ab926a40e1c7f8cb2d08267f/diff,workdir=/var/lib/docker/overlay2/69df9ec3d077f103f8a9e1eab68cc60c18dcae03ab926a40e1c7f8cb2d08267f/work
	RootMountInfoVFSOptionsContainsDocker = RootMountInfoVFSOptionsContains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "rootfs vfs otions",
			Info:   "rootfs's mountinfo vfs options contains 'docker', e.g., /var/lib/docker/overlay2...",
			ExeEnv: exeenv.InContainer,
		},
		Expected: "docker",
	}
)

type HostsMountInfoRootContains struct {
	prerequisite.BasePrerequisite
	Expected string
}

func (p *HostsMountInfoRootContains) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	info, err := mountinfo.HostsMount()
	if err != nil {
		return false, fmt.Errorf("failed to check %s caused by getting mountinfo of /etc/hosts: %w", p.Name, err)
	}
	p.Satisfied = strings.Contains(info.Root, p.Expected)
	p.Checked = true
	return p.Satisfied, nil
}

var (
	// HostsMountInfoRootContainsDocker
	//https://github.com/moby/moby/blob/v28.4.0/daemon/container_operations_unix.go#L552
	//https://github.com/moby/moby/blob/v28.4.0/daemon/config/config_linux.go#L189
	//https://github.com/moby/moby/blob/v28.4.0/container/container.go#L407
	HostsMountInfoRootContainsDocker = HostsMountInfoRootContains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "/etc/hosts",
			Info:   "/etc/hosts's mountinfo root contains 'docker', e.g., 814 696 259:2 /var/lib/docker/containers/44bec6602ccfae7458bcd71279beafc287d9fde509fa861377e162270d4cd92f/hosts /etc/hosts rw,relatime - ext4 /dev/nvme0n1p2 rw,errors=remount-ro",
			ExeEnv: exeenv.InContainer,
		},
		Expected: "docker",
	}
)
