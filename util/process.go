package util

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func GetProcessNameByPid(pid int) (processName string, err error) {
	path, err := GetProcessPathByPid(pid)
	if err != nil {
		return
	}
	processName = filepath.Base(path)
	return
}

func GetProcessPathByPid(pid int) (path string, err error) {
	return getProcessPath(strconv.Itoa(pid))
}

func getProcessPath(pid string) (path string, err error) {
	exe := fmt.Sprintf("/proc/%s/exe", pid)
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

func getSelfPid() (pid int, err error) {
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
	selfPid, err := getSelfPid()
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
	cmdline, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	comm, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	if strings.HasPrefix(string(cmdline), string(comm)) {
		if strings.HasPrefix(string(cmdline), "/bin/sh") {
			shebang = true
			return
		}
	}
	return
}
