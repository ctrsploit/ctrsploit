package docker_compose

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	docker_compose "github.com/ctrsploit/ctrsploit/pkg/version/docker-compose"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type VersionBetween struct {
	prerequisite.BasePrerequisite
	Min string
	Max string
}

func (p *VersionBetween) Check() (bool, error) {
	return p.CheckTemplate(func() {
		version, err := docker_compose.Version()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("getting docker compose version: %w", err))
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
	VulnerableToCVE_2025_62725 = &VersionBetween{
		Min: "2.0.0",
		Max: "2.40.1",
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "Docker Compose Vulnerable to CVE-2025-62725",
			Info:   "Docker Compose >= v2.0.0 and <= v2.40.1 is vulnerable to CVE-2025-62725",
			ExeEnv: exeenv.Local,
		},
	}
)
