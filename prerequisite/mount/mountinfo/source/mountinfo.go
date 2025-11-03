package source

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
	return p.CheckTemplate(func() {
		info, err := mountinfo.RootMount()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("getting root's mountinfo: %w", err))
			return
		}
		p.Satisfied = strings.Contains(info.Source, p.Expected)
		return
	})
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
		return false, fmt.Errorf("failed to check [%s], caused by getting root's mountinfo: %w", p.GetName(), err)
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
	// RootMountInfoVFSOptionsContainsContainerd
	//https://github.com/containerd/containerd/blob/v2.1.4/defaults/defaults_unix.go#L26
	RootMountInfoVFSOptionsContainsContainerd = RootMountInfoVFSOptionsContains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "rootfs vfs otions",
			Info:   "rootfs's mountinfo vfs options contains 'containerd', e.g., lowerdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/...",
			ExeEnv: exeenv.InContainer,
		},
		Expected: "containerd",
	}
)
