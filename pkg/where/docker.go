package where

import (
	"bytes"
	"ctrsploit/pkg/graphdriver"
	"ctrsploit/pkg/mountinfo"
	"ctrsploit/util"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const (
	PatternDockerHostname = "^[0-9a-f]{12}$"
)

type Docker struct {
	DockerEnvFileExists            bool
	GraphDriver                    graphdriver.GraphDriver
	RootfsContainsDocker           bool
	CgroupContainsDocker           bool
	HostsMountSourceContainsDocker bool
	HostnameMatchPattern           bool
}

func (d *Docker) CheckDockerEnvExists() {
	d.DockerEnvFileExists = util.CheckFileExists("/.dockerenv")
}

// CheckMountInfo rootfs contains "docker"
func (d *Docker) CheckMountInfo() (err error) {
	err = d.GraphDriver.Init()
	if err != nil {
		return
	}
	d.RootfsContainsDocker = strings.Contains(d.GraphDriver.Rootfs, "docker")
	return
}

// CheckCgroup Only works in cgroup v1
func (d *Docker) CheckCgroup() (err error) {
	content, err := ioutil.ReadFile("/proc/self/cgroup")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	d.CgroupContainsDocker = bytes.Contains(content, []byte("docker"))
	return
}

func (d *Docker) CheckHostsMountSourceContainsDocker() (err error) {
	mount, err := mountinfo.HostsMount()
	if err != nil {
		return
	}
	d.HostsMountSourceContainsDocker = strings.Contains(mount.Root, "docker")
	return
}

func (d *Docker) CheckHostnameMatchPattern() (err error) {
	hostname, err := os.Hostname()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	d.HostnameMatchPattern, err = regexp.MatchString(PatternDockerHostname, hostname)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}

func (d *Docker) IsIn() (in bool, err error) {
	d.CheckDockerEnvExists()
	if d.DockerEnvFileExists {
		in = true
	}

	err = d.CheckMountInfo()
	if err != nil {
		return
	}
	if d.RootfsContainsDocker {
		in = true
	}

	err = d.CheckCgroup()
	if err != nil {
		return
	}
	if d.CgroupContainsDocker {
		in = true
	}

	err = d.CheckHostsMountSourceContainsDocker()
	if err != nil {
		return
	}
	if d.HostsMountSourceContainsDocker {
		in = true
	}

	err = d.CheckHostnameMatchPattern()
	if err != nil {
		return
	}
	return
}

func (d *Docker) Init() (err error) {
	_, err = d.IsIn()
	return
}
