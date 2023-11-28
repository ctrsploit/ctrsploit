package capability

import (
	"github.com/containerd/containerd/pkg/cap"
	"github.com/ctrsploit/ctrsploit/pkg/capability"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/slice"
)

type Capability struct {
	ExpectedCapability string
	prerequisite.BasePrerequisite
}

var ContainsCapSysAdmin = Capability{
	ExpectedCapability: "CAP_SYS_ADMIN",
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name: "CAP_SYS_ADMIN",
		Info: "Container with cap_sys_admin is dangerous",
	},
}

func (p *Capability) Check() (err error) {
	err = p.BasePrerequisite.Check()
	if err != nil {
		return
	}
	caps, err := capability.GetPid1Capability()
	if err != nil {
		return
	}
	capsParsed, _ := cap.FromBitmap(caps)
	p.Satisfied = slice.In(p.ExpectedCapability, capsParsed)
	return
}
