package pid_host

import "syscall"

func getPC(reg syscall.PtraceRegs) uint64 {
	return reg.Pc
}
