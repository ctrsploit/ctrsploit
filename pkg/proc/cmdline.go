package proc

import (
	"fmt"
	"os"
	"strings"
)

func Cmdline() ([]string, error) {
	content, err := os.ReadFile("/proc/self/cmdline")
	if err != nil {
		return nil, fmt.Errorf("error reading /proc/self/cmdline: %v", err)
	}
	return strings.Split(string(content), "\x00"), nil
}
