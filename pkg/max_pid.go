package pkg

import (
	"github.com/ctrsploit/ctrsploit/internal"
)

func MaxPid() (int, error) {
	return internal.ReadIntFromFile("/proc/sys/kernel/pid_max")
}
