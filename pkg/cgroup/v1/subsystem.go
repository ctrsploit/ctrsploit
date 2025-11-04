package v1

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

/*
ParseCgroupFile is copied from github.com/containerd/cgroups/utils.go:parseCgroupFromReaderUnified(), but without splitting ","

because cgroup subsystem cannot be mounted after split by ','

https://github.com/torvalds/linux/blob/v5.4/kernel/cgroup/cgroup-v1.c#L1186
https://github.com/torvalds/linux/blob/v5.4/kernel/cgroup/cgroup-v1.c#L1141

e.g.: net_prio cannot be mounted in the non-init cgroup ns

root@cve-2022-0492:~# mkdir /tmp/net_prio
root@cve-2022-0492:~# unshare -UrCm
root@cve-2022-0492:~# mount -t cgroup -o net_prio none /tmp/net_prio
mount: /tmp/net_prio: permission denied.
root@cve-2022-0492:~# bpftrace -e 'kretprobe:cgroup1_get_tree /comm=="mount"/ {printf("%s:%d", comm, retval);printf("%s\n", kstack); }'
Attaching 1 probe...
mount:-1

	kretprobe_trampoline+0
	do_mount+1969
	ksys_mount+130
	__x64_sys_mount+37
	do_syscall_64+87
	entry_SYSCALL_64_after_hwframe+68

cgroup1_root_to_use() iterate across each hierarchy, net_prio is not a hierarchy:

root@cve-2022-0492:~# cat /proc/self/cgroup
12:cpu,cpuacct:/
11:perf_event:/
10:blkio:/
9:cpuset:/
8:pids:/
7:hugetlb:/
6:memory:/
5:freezer:/
4:rdma:/
3:net_cls,net_prio:/
2:devices:/
1:name=systemd:/
0::/
*/
func ParseCgroupFile(r io.Reader) (map[string]string, error) {
	var (
		subsystems = make(map[string]string)
		s          = bufio.NewScanner(r)
	)
	for s.Scan() {
		var (
			text  = s.Text()
			parts = strings.SplitN(text, ":", 3)
		)
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid cgroup entry: %q", text)
		}
		//for _, subs := range strings.Split(parts[1], ",") {
		//	if subs == "" {
		//		unified = parts[2]
		//	} else {
		//		subsystems[subs] = parts[2]
		//	}
		//}
		subsystems[parts[1]] = parts[2]
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return subsystems, nil
}

func (c CgroupV1) ListSubsystemsQuick(procCgroupPath string) (subsystems map[string]string, err error) {
	f, err := os.Open(procCgroupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", procCgroupPath, err)
	}
	defer f.Close()
	subsystems, err = ParseCgroupFile(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", procCgroupPath, err)
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

/*
use /proc/1/cgroup instead of /proc/self/cgroup
because when execute in host, /proc/self/cgroup in sub cgroup ns shows top level sub system, but actually not.
and in host, we can use /proc/1/cgroup to shows real levels.

root@cve-2022-0492:~# unshare -UrCm
root@cve-2022-0492:~# cat /proc/1/cgroup
12:perf_event:/
11:hugetlb:/
10:net_cls,net_prio:/
9:rdma:/
8:pids:/../../..
7:devices:/..
6:cpu,cpuacct:/..
5:memory:/../../..
4:blkio:/..
3:cpuset:/
2:freezer:/
1:name=systemd:/../../../init.scope
0::/../../../init.scope
root@cve-2022-0492:~# cat /proc/self/cgroup
12:perf_event:/
11:hugetlb:/
10:net_cls,net_prio:/
9:rdma:/
8:pids:/
7:devices:/
6:cpu,cpuacct:/
5:memory:/
4:blkio:/
3:cpuset:/
2:freezer:/
1:name=systemd:/
0::/
*/
func listTopLevelSubSystemQuick() (topLevelSubSystems []string, err error) {
	var c CgroupV1
	subsystemsSupport, err := c.ListSubsystemsQuick("/proc/1/cgroup")
	if err != nil {
		return nil, fmt.Errorf("ListSubsystemsQuick() failed: %w", err)
	}
	for name, path := range subsystemsSupport {
		if c.IsTopQuick(path) {
			topLevelSubSystems = append(topLevelSubSystems, name)
		}
	}
	return
}

func ListTopLevelSubSystem() (topLevelSubSystems []string, err error) {
	var c CgroupV1
	subsystems, err := listTopLevelSubSystemQuick()
	if err != nil {
		return nil, fmt.Errorf("listTopLevelSubSystemQuick: failed to list sub systems: %w", err)
	}
	for _, subsystemName := range subsystems {
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
