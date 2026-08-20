package ptraceinject

import "syscall"

func getPC(reg syscall.PtraceRegs) uint64 {
	return reg.Rip
}
