package buildkitd

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ctrsploit/ctrsploit/pkg/version/buildkitd"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type VersionBetween struct {
	prerequisite.BasePrerequisite
	Addr string
	Min  string
	Max  string
}

func (p *VersionBetween) Check() (bool, error) {
	return p.CheckTemplate(func() {
		version, err := buildkitd.Version(p.Addr)
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("getting buildkitd version: %w", err))
			return
		}
		rule := fmt.Sprintf(">= %s, <= %s", p.Min, p.Max)
		constraint, err := semver.NewConstraint(rule)
		if err != nil {
			err = fmt.Errorf("failed to parse constraint %s: %w", rule, err)
			// it's fatal if the constraint is invalid, it means the code is wrong
			awesome_error.CheckFatal(err)
		}
		p.Satisfied = constraint.Check(version)
	})
}

//goland:noinspection GoSnakeCaseUsage
var (
	VulnerableToCVE_2024_23650 = func(addr string) *VersionBetween {
		return &VersionBetween{
			Addr: addr,
			Min:  "0.0.0",
			Max:  "0.12.4",
			BasePrerequisite: prerequisite.BasePrerequisite{
				Name:   "BuildKitd Vulnerable to CVE-2024-23650",
				Info:   "Buildkitd <= v0.12.4 is vulnerable to CVE-2024-23650",
				ExeEnv: exeenv.Remote,
			},
		}
	}
)
