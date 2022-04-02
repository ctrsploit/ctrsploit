package util

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
	"path/filepath"
	"strconv"
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
