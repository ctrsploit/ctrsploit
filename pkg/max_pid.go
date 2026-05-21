package pkg

import (
	"github.com/ctrsploit/ctrsploit/pkg/fileutil"
)

func MaxPid() (int, error) {
	return fileutil.ReadIntFromFile("/proc/sys/kernel/pid_max")
}
