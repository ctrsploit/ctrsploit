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
	HostnameMatchPattern           bool
}

func NewDocker() *Docker {
	d := &Docker{}
	d.CheckAll()
	return d
}

func (d *Docker) Is() (bool, error) {
	return d.DockerEnvFileExists ||
		d.RootfsContainsDocker ||
		d.CgroupContainsDocker ||
		d.HostsMountSourceContainsDocker ||
		d.HostnameMatchPattern, nil
}

func (d *Docker) CheckAll() {
	err := d.checkDockerEnvExists()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	err = d.checkMountInfo()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	err = d.checkCgroup()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	err = d.checkHostsMountSourceContainsDocker()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	err = d.checkHostnameMatchPattern()
	if err != nil {
		awesome_error.CheckWarning(err)
	}
	return
}

func (d *Docker) checkDockerEnvExists() (err error) {
	d.DockerEnvFileExists, err = internal.CheckPathExists("/.dockerenv")
	if err != nil {
		return fmt.Errorf("error on checking /.dockerenv exists: %w", err)
	}
	return
}

// checkMountInfo rootfs contains "docker"
func (d *Docker) checkMountInfo() (err error) {
	info, err := mountinfo.RootMount()
	if err != nil {
		return fmt.Errorf("error getting root's mount info: %w", err)
	}
	// device mapper: /dev/mapper/docker-253:0-...
	if strings.Contains(info.Source, "docker") ||
		// overlay: lowerdir=/var/lib/docker/overlay2...
		strings.Contains(info.VFSOptions, "docker") {
		d.RootfsContainsDocker = true
	}
	return
}

// checkCgroup Only works in cgroup v1
func (d *Docker) checkCgroup() (err error) {
	content, err := os.ReadFile("/proc/self/cgroup")
	if err != nil {
		return fmt.Errorf("error reading /proc/self/cgroup: %w", err)
	}
	d.CgroupContainsDocker = bytes.Contains(content, []byte("docker"))
	return
}

func (d *Docker) checkHostsMountSourceContainsDocker() (err error) {
	mount, err := mountinfo.HostsMount()
	if err != nil {
		return fmt.Errorf("error getting mountinfo of /etc/hosts: %w", err)
	}
	d.HostsMountSourceContainsDocker = strings.Contains(mount.Root, "docker")
	return
}

func (d *Docker) checkHostnameMatchPattern() (err error) {
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("error getting hostname: %w", err)
	}
	d.HostnameMatchPattern, err = regexp.MatchString(PatternDockerHostname, hostname)
	if err != nil {
		return fmt.Errorf("error checking hostname match pattern: %w", err)
	}
	return
}
