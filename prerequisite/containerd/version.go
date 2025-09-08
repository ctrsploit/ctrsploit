package containerd

import (
	"github.com/ctrsploit/ctrsploit/pkg/version/containerd"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type VersionEqualTo struct {
	ExpectedVersion string
	prerequisite.BasePrerequisite
}

var VersionEqualToV2_1_0 = VersionEqualTo{
	ExpectedVersion: "v2.1.0",
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name: "containerd==v2.1.0",
		Info: "Containerd version must be 2.1.0",
	},
}

func (p *VersionEqualTo) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	v, err := containerd.GetVersionBySock()
	if err != nil {
		return false, err
	}
	p.Satisfied = v.Version == p.ExpectedVersion
	p.Checked = true
	return p.Satisfied, nil
}
