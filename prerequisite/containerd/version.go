package containerd

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ctrsploit/ctrsploit/pkg/version/containerd"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type VersionEqualTo struct {
	ExpectedVersion string
	prerequisite.BasePrerequisite
}

//goland:noinspection GoSnakeCaseUsage
var VersionEqualToV2_1_0 = VersionEqualTo{
	ExpectedVersion: "v2.1.0",
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name: "containerd==v2.1.0",
		Info: "Containerd version must be 2.1.0",
	},
}

func (p *VersionEqualTo) Check() (bool, error) {
	return p.CheckTemplate(func() {
		v, err := containerd.GetVersionBySock()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("getting version from containerd: %v", err))
			return
		}
		p.Satisfied = v.Version == p.ExpectedVersion
		return
	})
}

// Version checks if containerd version matches a semver constraint
type Version struct {
	prerequisite.BasePrerequisite
	Constraint *semver.Constraints
}

func (p *Version) Check() (bool, error) {
	return p.CheckTemplate(func() {
		v, err := containerd.GetVersionBySock()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("getting version from containerd: %v", err))
			return
		}
		ver, err := semver.NewVersion(v.Version)
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("parsing containerd version %s: %v", v.Version, err))
			return
		}
		p.Satisfied = p.Constraint.Check(ver)
		return
	})
}

// VulnerableTo40635 checks if containerd is vulnerable to CVE-2024-40635
// Vulnerable versions: < 1.6.38, < 1.7.27, < 2.0.4
//
//goland:noinspection GoSnakeCaseUsage
var VulnerableTo40635 = func() *Version {
	cons, err := semver.NewConstraint("(>= 1.6.0, < 1.6.38) || (>= 1.7.0, < 1.7.27) || (>= 2.0.0, < 2.0.4)")
	if err != nil {
		panic(err)
	}
	return &Version{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "containerd vulnerable to CVE-2024-40635",
			Info:   "containerd < 1.6.38, < 1.7.27, or < 2.0.4 is vulnerable to CVE-2024-40635 (Integer overflow in User ID handling)",
			ExeEnv: exeenv.InHost,
		},
		Constraint: cons,
	}
}()
