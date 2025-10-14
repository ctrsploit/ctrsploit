package runc

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ctrsploit/ctrsploit/pkg/version/runc"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type Version struct {
	prerequisite.BasePrerequisite
	Constraint *semver.Constraints
}

func (p *Version) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	p.Checked = true
	ver, err := runc.GetVersionByCliVersion()
	if err != nil {
		return false, fmt.Errorf("could not get version from cli: %w", err)
	}
	var errs []error
	p.Satisfied, errs = p.Constraint.Validate(ver)
	if len(errs) > 0 {
		return false, fmt.Errorf("constraint validation errors: %v", errs)
	}
	return p.Satisfied, nil
}

var (
	// NotVulnerableTo5736ByVersion is a sufficient condition
	NotVulnerableTo5736ByVersion = func() *Version {
		cons, err := semver.NewConstraint(">= v1.0.0-rc7")
		if err != nil {
			panic(err)
		}
		return &Version{
			BasePrerequisite: prerequisite.BasePrerequisite{
				Name: "runc not vulnerable to CVE-2019-5736",
				// https://github.com/opencontainers/runc/commit/0a8e4117e7f715d5fbeef398405813ce8e88558b
				Info:   "runc >= v1.0.0-rc7 is not vulnerable to CVE-2019-5736",
				ExeEnv: exeenv.InHost,
			},
			Constraint: cons,
		}
	}
	// MayBeVulnerableTo5736ByVersion is a necessary condition but not sufficient condition,
	// because downstream vendor may backport the patch to an older version.
	MayBeVulnerableTo5736ByVersion = prerequisite.Not(NotVulnerableTo5736ByVersion())
)
