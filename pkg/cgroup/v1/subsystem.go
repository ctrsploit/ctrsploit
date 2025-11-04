package v1

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/cgroups"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const releaseAgent = "release_agent"
const DefaultMountPoint = "/sys/fs/cgroup"

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
	fileInfo, err := os.ReadDir(mountpoint)
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

func (c CgroupV1) ListSubsystemsDeprecated(procCgroupPath string) (subsystems map[string]string, err error) {
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

/*
IsTopQuick fails in the sub cgroup ns
root@73313224d1ef:/# unshare -UrCm /bin/bash
root@73313224d1ef:/# cat /proc/1/cgroup
12:freezer:/
11:hugetlb:/
10:memory:/
9:cpu,cpuacct:/
8:pids:/
7:cpuset:/
6:blkio:/
5:perf_event:/
4:net_cls,net_prio:/
3:rdma:/
2:devices:/
1:name=systemd:/
0::/
*/
func (c CgroupV1) IsTopQuick(subsystemPath string) (top bool) {
	return subsystemPath == "/"
}

func ListTopLevelSubSystem() (topLevelSubSystems []string, err error) {
	var c CgroupV1
	subsystemsSupport, err := c.ListSubsystems(DefaultMountPoint)
	if err != nil {
		return nil, fmt.Errorf("ListTopLevelSubSystem: failed to list sub systems: %w", err)
	}
	for _, subsystemName := range subsystemsSupport {
		// add this to be more accurate on the host
		if !c.IsTopQuick(subsystemName) {
			continue
		}
		is, err := c.IsTop(DefaultMountPoint, subsystemName)
		if err != nil {
			return nil, fmt.Errorf("ListTopLevelSubSystem: failed to list sub system %s: %w", subsystemName, err)
		}
		if is {
			topLevelSubSystems = append(topLevelSubSystems, subsystemName)
		}
	}
	return
}
