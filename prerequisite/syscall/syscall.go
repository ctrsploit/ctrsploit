package syscall

import (
	"errors"
	"syscall"

	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type HasPerm struct {
	prerequisite.BasePrerequisite
	SyscallNumber int
}

func (p *HasPerm) Check() (bool, error) {
	return p.CheckTemplate(func() (bool, error) {
		const invalidArgument = ^uintptr(0)
		_, _, errno := syscall.Syscall(uintptr(p.SyscallNumber), invalidArgument, 0, 0)
		// TODO: if syscall unsupported by kernel?
		p.Satisfied = !errors.Is(errno, syscall.EPERM)
		return p.Satisfied, p.Err
	})
}

var (
	Unshare = HasPerm{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "unshare",
			Info:   "available to call unshare syscall",
			ExeEnv: exeenv.InContainer,
		},
		SyscallNumber: syscall.SYS_UNSHARE,
	}
)
