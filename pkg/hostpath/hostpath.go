package hostpath

import (
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver"
	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"strings"
)

type Path struct {
	ContainerPath string
	HostPath      string
	Type          int
}

const (
	TypeUnknown = iota
	TypeRootfs
	TypeNetworkFiles
	TypeUserCustomBindMount
)

func RootFs() (path string, err error) {
	g := graphdriver.GraphDriver{}
	err = g.Init()
	if err == nil {
		path = g.Rootfs
	}
	return
}

func WritableAccessible() (paths []Path, err error) {
	// rootfs: get the host path of the container rootfs
	rootfs, err := RootFs()
	if err == nil {
		paths = append(paths, Path{
			ContainerPath: "/",
			HostPath:      rootfs,
			Type:          TypeRootfs,
		})
	}
	infos, err := mountinfo.MountInfo()
	if err != nil {
		return paths, err
	}
	for _, info := range infos {
		if strings.Contains(info.Options, "ro") {
			continue
		}
		switch info.Mountpoint {
		// handle network files
		case "/etc/hosts", "/etc/hostname", "/etc/resolv.conf":
			paths = append(paths, Path{
				ContainerPath: info.Mountpoint,
				HostPath:      info.Root,
				Type:          TypeNetworkFiles,
			})
		case "/",
			"/proc", "/proc/bus", "/proc/fs", "/proc/irq", "/proc/sys", "/proc/sysrq-trigger", "/proc/asound",
			"/proc/acpi", "/proc/kcore", "/proc/keys", "/proc/timer_list", "/proc/sched_debug",
			"/sys", "/sys/firmware", "/sys/devices/virtual/powercap",
			"/dev", "/dev/pts", "/dev/mqueue", "/dev/shm", "/dev/console":
			continue
		default:
			// Append the custom bind mount.
			paths = append(paths, Path{
				ContainerPath: info.Mountpoint,
				HostPath:      info.Root,
				Type:          TypeUserCustomBindMount,
			})
		}
	}
	return
}
