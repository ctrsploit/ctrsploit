package seccomp

import (
	"github.com/ctrsploit/ctrsploit/pkg/version"
	"golang.org/x/sys/unix"
	"syscall"
)

type Syscall struct {
	Number           int
	KernelMinVersion string
	KernelMaxVersion string
	DockerMinVersion version.Docker
	DockerMaxVersion version.Docker
}

func (s Syscall) Enabled() bool {
	_, _, errno := syscall.RawSyscall(
		uintptr(s.Number),
		0,
		0,
		0,
	)
	// return errno == unix.EFAULT
	return errno != unix.EPERM
}

var (
	// IOURingSetup
	// kernel enable: https://elixir.bootlin.com/linux/v5.1-rc1/source/include/linux/syscalls.h#L318
	// seccomp disable: https://github.com/moby/moby/pull/46762
	// seccomp enable: https://github.com/moby/moby/pull/39415
	IOURingSetup = Syscall{
		Number:           unix.SYS_IO_URING_SETUP,
		KernelMinVersion: "v5.1-rc1",
		DockerMaxVersion: version.NewDocker("25.0.0-beta.1"),
	}
)
