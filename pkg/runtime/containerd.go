package runtime

import (
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type Containerd struct {
	RootfsContainsContainerd              bool
	HostsMountSourceContainsNerdctl       bool
	HostnameMountSourceContainsContainerd bool
	ProcNetUnixContainsContainerdSock     bool
}

func NewContainerd() *Containerd {
	r := &Containerd{}
	r.CheckAll()
	return r
}

func (r *Containerd) Is() (bool, error) {
	return r.RootfsContainsContainerd ||
			r.HostsMountSourceContainsNerdctl ||
			r.HostnameMountSourceContainsContainerd ||
			r.ProcNetUnixContainsContainerdSock,
		nil
}

func (r *Containerd) CheckAll() {
	err := r.checkRootfsMountInfo()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	err = r.checkHostsMountSourceContainsNerdctl()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	err = r.checkHostnameMountSourceContainsContainerd()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	err = r.checkProcUnixNetContainsContainerdSock()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
}

// checkRootfsMountInfo rootfs contains "containerd"
func (r *Containerd) checkRootfsMountInfo() (err error) {
	info, err := mountinfo.RootMount()
	if err != nil {
		return fmt.Errorf("error getting root's mount info: %w", err)
	}
	// overlay: lowerdir=/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs/snapshots/...
	if strings.Contains(info.VFSOptions, "containerd") {
		r.RootfsContainsContainerd = true
	}
	return
}

func (r *Containerd) checkHostsMountSourceContainsNerdctl() (err error) {
	mount, err := mountinfo.HostsMount()
	if err != nil {
		return fmt.Errorf("error getting mountinfo of /etc/hosts: %w", err)
	}
	// note: the host path of mount source is mount.Root instead of mount.Source
	r.HostsMountSourceContainsNerdctl = strings.Contains(mount.Root, "nerdctl")
	return
}

func (r *Containerd) checkHostnameMountSourceContainsContainerd() (err error) {
	mount, err := mountinfo.Hostname()
	if err != nil {
		return fmt.Errorf("error getting mountinfo of /etc/hostname: %w", err)
	}
	// note: the host path of mount source is mount.Root instead of mount.Source
	if strings.Contains(mount.Root, "containerd") ||
		strings.Contains(mount.Root, "nerdctl") {
		r.HostnameMountSourceContainsContainerd = true
	}
	return
}

func (r *Containerd) checkProcUnixNetContainsContainerdSock() (err error) {
	content, err := os.ReadFile("/proc/net/unix")
	if err != nil {
		return fmt.Errorf("error reading /proc/net/unix: %w", err)
	}
	if strings.Contains(string(content), "containerd.sock") &&
		!strings.Contains(string(content), "docker.sock") {
		r.ProcNetUnixContainsContainerdSock = true
	}
	return
}
