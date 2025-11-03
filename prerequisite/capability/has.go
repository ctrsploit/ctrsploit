package capability

import (
	"fmt"

	"github.com/containerd/containerd/pkg/cap"
	"github.com/ctrsploit/ctrsploit/pkg/capability"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/slice"
)

type Has struct {
	ExpectedCapability string
	Pid                []string
	// check CapBnd or CapEff, CapBnd for check vul exists, CapEff for check vul exploitable
	CapType cap.Type
	prerequisite.BasePrerequisite
}

func HasBnd(name string) Has {
	return Has{
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

func HasEff(name string) Has {
	return Has{
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
	CapSysAdminBnd      = HasBnd("CAP_SYS_ADMIN")
	CapSysAdminEff      = HasEff("CAP_SYS_ADMIN")
	CapDacReadSearchBnd = HasBnd("CAP_DAC_READ_SEARCH")
	CapDacReadSearchEff = HasEff("CAP_DAC_READ_SEARCH")
	CapSysPtraceBnd     = HasBnd("CAP_SYS_PTRACE")
	CapSysPtraceEff     = HasEff("CAP_SYS_PTRACE")
	CapBpfBnd           = HasBnd("CAP_BPF")
	// CapBpfEff CAP_BPF: CAP_BPF load ebpf program
	CapBpfEff     = HasEff("CAP_BPF")
	CapPerfmonBnd = HasBnd("CAP_PERFMON")
	// CapPerfmonEff CAP_PERFMON: attach to kprobes, uprobes, tracepoints
	CapPerfmonEff     = HasEff("CAP_PERFMON")
	CapNetRawBnd      = HasBnd("CAP_NET_RAW")
	CapNetRawEff      = HasEff("CAP_NET_RAW")
	CapSysResourceBnd = HasBnd("CAP_SYS_RESOURCE")
	CapSysResourceEff = HasEff("CAP_SYS_RESOURCE")
)

func (p *Has) Check() (bool, error) {
	return p.CheckTemplate(func() {
		for _, pid := range p.Pid {
			caps, err := capability.GetCapabilityByPid(pid, p.CapType)
			if err != nil {
				p.Err = p.WrapErr(fmt.Errorf("getting capability of %s: %w", pid, err))
				return
			}
			capsParsed, _ := cap.FromBitmap(caps)
			if slice.In(p.ExpectedCapability, capsParsed) {
				p.Satisfied = true
				break
			}
		}
		return
	})
}
