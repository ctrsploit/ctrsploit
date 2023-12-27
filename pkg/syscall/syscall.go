package syscall

import (
	"errors"
	"github.com/ctrsploit/ctrsploit/pkg/version/version"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"golang.org/x/sys/unix"
	"syscall"
)

type Syscall struct {
	Number  int
	Version map[version.TypeState][]map[version.TypeSoftware][]version.Version
}

func (s Syscall) State() (state version.TypeState) {
	_, _, errno := syscall.RawSyscall(uintptr(s.Number), 0, 0, 0)
	log.Logger.Debugf("syscall %d errno: %+v", s.Number, errno)
	switch {
	case errors.Is(errno, unix.EPERM):
		state = version.StateDisable
	case errors.Is(errno, unix.EFAULT):
		state = version.StateValid
	case errors.Is(errno, unix.ENOSYS):
		state = version.StateUnsupported
	default:
		state = version.StateUnknown
	}
	return
}

func (s Syscall) Enabled() bool {
	return s.State() == version.StateValid
}

func (s Syscall) RangeOfSoftware(software version.TypeSoftware) (versions []version.Version) {
	// state := s.State()
	// for _, vers := range s.Version[state] {}
	return nil
}

func (s Syscall) Range() (result []map[version.TypeSoftware][]version.Version) {
	return
}
