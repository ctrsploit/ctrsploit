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

func (vb *VersionBetween) Check() (bool, error) {
	if vb.Checked {
		return vb.Satisfied, nil
	}
	version, err := buildkitd.Version(vb.Addr)
	if err != nil {
		return false, err
	}
	constraint, err := semver.NewConstraint(fmt.Sprintf(">= %s, <= %s", vb.Min, vb.Max))
	if err != nil {
		awesome_error.CheckErr(err)
		return false, err
	}
	var e []error
	vb.Satisfied, e = constraint.Validate(version)
	if len(e) > 0 {
		err = fmt.Errorf("failed to validate version %s: %v", version.String(), e)
		awesome_error.CheckErr(err)
		return false, err
	}
	vb.Checked = true
	return vb.Satisfied, nil
}

var VulnerableToCVE_2024_23650 = func(addr string) *VersionBetween {
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
