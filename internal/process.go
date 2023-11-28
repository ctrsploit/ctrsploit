package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func GetProcessNameByPid(pid int) (processName string, err error) {
	path, err := GetProcessPath(pid)
	if err != nil {
		return
	}
	processName = filepath.Base(path)
	return
}

func GetProcessPath(pid int) (path string, err error) {
	exe := fmt.Sprintf("/proc/%d/exe", pid)
	path, err = filepath.EvalSymlinks(exe)
	if err != nil {
		if errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrNotExist) {
			return
		} else {
			awesome_error.CheckErr(err)
		}
		return
	}
	return
}

func GetSelfPid() (pid int, err error) {
	path, err := filepath.EvalSymlinks("/proc/self")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	process := strings.TrimPrefix(path, "/proc/")
	pid, err = strconv.Atoi(process)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}

func KillAll() (err error) {
	matches, err := filepath.Glob("/proc/[0-9]*")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	selfPid, err := GetSelfPid()
	if err != nil {
		return
	}
	for _, match := range matches {
		process := strings.TrimPrefix(match, "/proc/")
		pid, err := strconv.Atoi(process)
		if err != nil {
			awesome_error.CheckErr(err)
			continue
		}
		if pid == selfPid {
			continue
		}
		err = syscall.Kill(pid, syscall.Signal(9))
		if err != nil {
			awesome_error.CheckErr(err)
			continue
		}
	}
	return
}

func IsSheBang(pid int) (shebang bool, err error) {
	cmdline, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	args := bytes.Split(cmdline, []byte{0})
	lastArg := string(args[len(args)-2])
	comm, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	if strings.Contains(lastArg, strings.TrimSpace(string(comm))) {
		shebang = true
	}
	return
}
