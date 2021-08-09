package cgroup

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
)

const (
	unifiedMountpoint = "/sys/fs/cgroup"
)

//IsCgroupV2BorrowedFromRunc https://github.com/opencontainers/runc/blob/3f2f06dfe1b3289b01daa531964b4f0af49cdf2d/docs/cgroup-v2.md#am-i-using-cgroup-v2
func IsCgroupV2BorrowedFromRunc() bool {
	return cgroups.IsCgroup2UnifiedMode()
}

func getStatfs() (st unix.Statfs_t, err error) {
	err = unix.Statfs(unifiedMountpoint, &st)
	if err != nil {
		// TODO: userns not found?
		awesome_error.CheckErr(err)
		return
	}
	return
}

func isCgroupV1(st unix.Statfs_t) bool {
	log.Logger.Infof("st type: %x", st.Type)
	if st.Type == unix.CGROUP_SUPER_MAGIC || st.Type == unix.TMPFS_MAGIC {
		return true
	}
	return false
}

// IsCgroupV1
// borrowed from https://github.com/opencontainers/runc/blob/v1.0.0/libcontainer/cgroups/fs/fs.go#L74
func IsCgroupV1() bool {
	var st, pst unix.Stat_t

	// (1) it should be a directory...
	err := unix.Lstat(unifiedMountpoint, &st)
	if err != nil || st.Mode&unix.S_IFDIR == 0 {
		awesome_error.CheckWarning(err)
		return false
	}

	// (2) ... and a mount point ...
	err = unix.Lstat(filepath.Dir(unifiedMountpoint), &pst)
	if err != nil {
		awesome_error.CheckWarning(err)
		return false
	}

	if st.Dev == pst.Dev {
		// parent dir has the same dev -- not a mount point
		return false
	}

	// (3) ... of 'tmpfs' fs type.
	var fst unix.Statfs_t
	err = unix.Statfs(unifiedMountpoint, &fst)
	if err != nil || fst.Type != unix.TMPFS_MAGIC {
		awesome_error.CheckWarning(err)
		return false
	}

	// (4) it should have at least 1 entry ...
	dir, err := os.Open(unifiedMountpoint)
	if err != nil {
		awesome_error.CheckWarning(err)
		return false
	}
	names, err := dir.Readdirnames(1)
	if err != nil {
		awesome_error.CheckWarning(err)
		return false
	}
	if len(names) < 1 {
		return false
	}
	// ... which is a cgroup mount point.
	err = unix.Statfs(filepath.Join(unifiedMountpoint, names[0]), &fst)
	if err != nil || fst.Type != unix.CGROUP_SUPER_MAGIC {
		awesome_error.CheckWarning(err)
		return false
	}

	return true
}

func isCgroupV2(st unix.Statfs_t) bool {
	return st.Type == unix.CGROUP2_SUPER_MAGIC
}

func IsCgroupV2() (is bool) {
	st, err := getStatfs()
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	is = isCgroupV2(st)
	return
}
