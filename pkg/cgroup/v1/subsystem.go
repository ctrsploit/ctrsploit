package v1

import (
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const releaseAgent = "release_agent"

// IsTopOld
/*
borrowed from: https://www.kernel.org/doc/Documentation/cgroup-v1/cgroups.txt

 - release_agent: the path to use for release notifications (this file
   exists in the top cgroup only)
*/
func (c CgroupV1) IsTopOld(mountpoint, subsystemName string) (top bool, err error) {
	_, err = os.Lstat(mountpoint)
	if err != nil {
		return
	}
	_, err = os.Lstat(filepath.Join(mountpoint, subsystemName))
	if err != nil {
		return
	}
	path := filepath.Join(mountpoint, subsystemName, releaseAgent)
	_, err = os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		} else {
			return
		}
	} else {
		top = true
		return
	}
}

func (c CgroupV1) ListSubsystemsOld(mountpoint string) (subsystems []string, err error) {
	fileInfo, err := ioutil.ReadDir(mountpoint)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	for _, file := range fileInfo {
		if file.IsDir() {
			subsystems = append(subsystems, file.Name())
		}
	}
	return
}

func (c CgroupV1) ListSubsystems(procCgroupPath string) (subsystems map[string]string, err error) {
	subsystems, err = cgroups.ParseCgroupFile(procCgroupPath)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	for sub := range subsystems {
		if strings.HasPrefix(sub, "name=") {
			delete(subsystems, sub)
		}
		if sub == "" {
			delete(subsystems, sub)
		}
	}
	return
}

func (c CgroupV1) IsTop(subsystemPath string) (top bool) {
	return subsystemPath == "/"
}
