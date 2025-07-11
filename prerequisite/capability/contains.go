package capability

import (
	"fmt"
	"github.com/containerd/containerd/pkg/cap"
	"github.com/ctrsploit/ctrsploit/pkg/capability"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/slice"
)

type Contains struct {
	ExpectedCapability string
	Pid                []string
	// check CapBnd or CapEff, CapBnd for check vul exists, CapEff for check vul exploitable
	CapType cap.Type
	prerequisite.BasePrerequisite
}

func BndContainsCap(name string) Contains {
	return Contains{
		ExpectedCapability: name,
		Pid:                []string{"1", "self"},
		CapType:            cap.Bounding,
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   name,
			Info:   fmt.Sprintf("CapBnd has %s", name),
			ExeEnv: exeenv.InContainer,
		},
	}
}

func EffContainsCap(name string) Contains {
	return Contains{
		ExpectedCapability: name,
		Pid:                []string{"self"},
		CapType:            cap.Effective,
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   name,
			Info:   fmt.Sprintf("CapEff has %s", name),
			ExeEnv: exeenv.InContainer,
		},
	}
}

var (
	CapSysAdminBnd      = BndContainsCap("CAP_SYS_ADMIN")
	CapSysAdminEff      = EffContainsCap("CAP_SYS_ADMIN")
	CapDacReadSearchBnd = BndContainsCap("CAP_DAC_READ_SEARCH")
	CapDacReadSearchEff = EffContainsCap("CAP_DAC_READ_SEARCH")
)

func (p *Contains) Check() (err error) {
	err = p.BasePrerequisite.Check()
	if err != nil {
		return
	}
	for _, pid := range p.Pid {
		caps, err := capability.GetCapabilityByPid(pid, p.CapType)
		if err != nil {
			return err
		}
		capsParsed, _ := cap.FromBitmap(caps)
		if slice.In(p.ExpectedCapability, capsParsed) {
			p.Satisfied = true
			break
		}
	}
	return
}
