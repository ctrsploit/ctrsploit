package user

import (
	"fmt"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"golang.org/x/sys/unix"
)

// UIDEqualTo check current process's ruid/euid/suid equal to expected values
type UIDEqualTo struct {
	prerequisite.BasePrerequisite
	RUid *int // Optional, nil means not checking
	EUid *int // Optional, nil means not checking
	SUid *int // Optional, nil means not checking
}

func (p *UIDEqualTo) Check() (bool, error) {
	return p.CheckTemplate(func() {
		ruid, euid, suid := unix.Getresuid()

		checks := []struct {
			name     string
			expected *int
			actual   int
		}{
			{"ruid", p.RUid, ruid},
			{"euid", p.EUid, euid},
			{"suid", p.SUid, suid},
		}

		var mismatches []string
		for _, chk := range checks {
			if chk.expected != nil && *chk.expected != chk.actual {
				mismatches = append(mismatches, fmt.Sprintf("%s=%d (expected %d)", chk.name, chk.actual, *chk.expected))
			}
		}

		if len(mismatches) > 0 {
			p.Err = p.WrapErr(fmt.Errorf("uid mismatch: %s", strings.Join(mismatches, ", ")))
		} else {
			p.Satisfied = true
		}
		return
	})
}

type UIDOption func(*UIDEqualTo)

func NewUIDEqualToPrerequisite(name, info string, env int, opts ...UIDOption) UIDEqualTo {
	p := UIDEqualTo{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   name,
			Info:   info,
			ExeEnv: env,
		},
	}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

func WithRUid(uid int) UIDOption {
	return func(p *UIDEqualTo) {
		p.RUid = &uid
	}
}

func WithEUid(uid int) UIDOption {
	return func(p *UIDEqualTo) {
		p.EUid = &uid
	}
}

func WithSUid(uid int) UIDOption {
	return func(p *UIDEqualTo) {
		p.SUid = &uid
	}
}

var (
	EUid0 = NewUIDEqualToPrerequisite(
		"euid=0",
		"",
		exeenv.InContainer,
		WithEUid(0),
	)

	RUid0 = NewUIDEqualToPrerequisite(
		"ruid=0",
		"",
		exeenv.InContainer,
		WithRUid(0),
	)
)
