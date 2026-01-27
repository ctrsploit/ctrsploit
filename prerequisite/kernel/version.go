package kernel

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ctrsploit/ctrsploit/pkg/version/kernel"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type VersionConstraint struct {
	prerequisite.BasePrerequisite
	Constraint string
}

func (p *VersionConstraint) Check() (bool, error) {
	return p.CheckTemplate(func() {
		cons, err := semver.NewConstraint(p.Constraint)
		if err != nil {
			err = fmt.Errorf("failed to parse constraint %s: %w", p.Constraint, err)
			awesome_error.CheckFatal(err)
		}
		version, err := kernel.Version()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("get kernel version: %w", err))
			return
		}
		p.Satisfied = cons.Check(version)
	})
}
