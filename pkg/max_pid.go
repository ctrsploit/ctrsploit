package pkg

import "ctrsploit/util"

func MaxPid() (int, error) {
	return util.ReadIntFromFile("/proc/sys/kernel/pid_max")
}
