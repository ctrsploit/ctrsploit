package hostpath

import (
	"fmt"
	"strings"

	"github.com/ctrsploit/ctrsploit/pkg/hostpath/rootfs"
	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ctrsploit/ctrsploit/pkg/runtime"
	"github.com/ctrsploit/ctrsploit/pkg/storagedriver"
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
	runtimeType := runtime.GetType()
	storageDriverType, err := storagedriver.GetType()
	if err != nil {
		return "", fmt.Errorf("failed to get storage driver type: %w", err)
	}
	return rootfs.HostPath(runtimeType, storageDriverType)
}

func EtcHosts() (string, error) {
	info, err := mountinfo.GetMountByMountpoint("/etc/hosts")
	if err != nil {
		return "", fmt.Errorf("failed to get hosts info: %w", err)
	}
	return info.Root, nil
}

func WritableAccessible() (paths []Path, err error) {
	// rootfs: get the host path of the container rootfs
	rootfsHostPath, err := RootFs()
	if err == nil {
		paths = append(paths, Path{
			ContainerPath: "/",
			HostPath:      rootfsHostPath,
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
