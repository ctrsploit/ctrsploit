package v1

import (
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io/ioutil"
	"os"
	"path/filepath"
)

const releaseAgent = "release_agent"

// IsTop
/*
borrowed from: https://www.kernel.org/doc/Documentation/cgroup-v1/cgroups.txt

 - release_agent: the path to use for release notifications (this file
   exists in the top cgroup only)
*/
func (c CgroupV1) IsTop(mountpoint, subsystemName string) (top bool, err error) {
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

func (c CgroupV1) ListSubsystems(mountpoint string) (subsystems []string, err error) {
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
