package containerd

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/version/containerd"
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
