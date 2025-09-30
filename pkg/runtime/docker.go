package runtime

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const (
	PatternDockerHostname = "^[0-9a-f]{12}$"
)

type Docker struct {
	DockerEnvFileExists            bool
	RootfsContainsDocker           bool
	CgroupContainsDocker           bool
	HostsMountSourceContainsDocker bool
	ProcAttrCurrentContainsDocker  bool
	// both docker, nerdctl match this behavior, disable this method
	hostnameMatchPattern          bool
	ProcNetUnixContainsDockerSock bool
}

func NewDocker() *Docker {
	d := &Docker{}
	d.CheckAll()
	return d
}

func (r *Docker) Is() (bool, error) {
	return r.DockerEnvFileExists ||
			r.RootfsContainsDocker ||
			r.CgroupContainsDocker ||
			r.ProcAttrCurrentContainsDocker ||
			//r.hostnameMatchPattern || // both docker, nerdctl match this behavior
			r.HostsMountSourceContainsDocker ||
			r.ProcNetUnixContainsDockerSock,
		nil
}

func (r *Docker) CheckAll() {
	err := r.checkDockerEnvExists()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	err = r.checkRootfsMountInfo()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	err = r.checkCgroup()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	err = r.checkHostsMountSourceContainsDocker()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	err = r.checkProcAttrCurrentContainsDocker()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	err = r.checkProcAttrCurrentContainsDocker()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	err = r.checkHostnameMatchPattern()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
}

func (r *Docker) checkDockerEnvExists() (err error) {
	r.DockerEnvFileExists, err = internal.CheckPathExists("/.dockerenv")
	if err != nil {
		return fmt.Errorf("error on checking /.dockerenv exists: %w", err)
	}
	return
}

// checkRootfsMountInfo rootfs contains "docker"
func (r *Docker) checkRootfsMountInfo() (err error) {
	info, err := mountinfo.RootMount()
	if err != nil {
		return fmt.Errorf("error getting root's mount info: %w", err)
	}
	// device mapper: /dev/mapper/docker-253:0-...
	if strings.Contains(info.Source, "docker") ||
		// overlay: lowerdir=/var/lib/docker/overlay2...
		strings.Contains(info.VFSOptions, "docker") {
		r.RootfsContainsDocker = true
	}
	return
}

// checkCgroup Only works in cgroup v1
func (r *Docker) checkCgroup() (err error) {
	content, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return fmt.Errorf("error reading /proc/1/cgroup: %w", err)
	}
	r.CgroupContainsDocker = bytes.Contains(content, []byte("docker"))
	return
}

func (r *Docker) checkHostsMountSourceContainsDocker() (err error) {
	mount, err := mountinfo.HostsMount()
	if err != nil {
		return fmt.Errorf("error getting mountinfo of /etc/hosts: %w", err)
	}
	// note: the host path of mount source is mount.Root instead of mount.Source
	r.HostsMountSourceContainsDocker = strings.Contains(mount.Root, "docker")
	return
}

func (r *Docker) checkProcAttrCurrentContainsDocker() (err error) {
	content, err := os.ReadFile("/proc/1/attr/current")
	if err != nil {
		return fmt.Errorf("error reading /proc/1/attr/current: %w", err)
	}
	r.ProcAttrCurrentContainsDocker = strings.Contains(string(content), "docker")
	return
}

func (r *Docker) checkProcUnixNetContainsDockerSock() (err error) {
	content, err := os.ReadFile("/proc/net/unix")
	if err != nil {
		return fmt.Errorf("error reading /proc/net/unix: %w", err)
	}
	r.ProcNetUnixContainsDockerSock = strings.Contains(string(content), "docker.sock")
	return
}

func (r *Docker) checkHostnameMatchPattern() (err error) {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("error getting hostname: %w", err)
	}
	r.hostnameMatchPattern, err = regexp.MatchString(PatternDockerHostname, hostname)
	if err != nil {
		return fmt.Errorf("error checking hostname match pattern: %w", err)
	}
	return
}
