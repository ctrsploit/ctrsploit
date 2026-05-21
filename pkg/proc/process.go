package proc

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func IsSheBang(pid int) (bool, error) {
	cmdline, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return false, err
	}
	args := bytes.Split(cmdline, []byte{0})
	if len(args) < 2 {
		return false, nil
	}
	lastArg := string(args[len(args)-2])
	comm, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return false, err
	}
	return strings.Contains(lastArg, strings.TrimSpace(string(comm))), nil
}
