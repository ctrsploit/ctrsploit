package kubelet

import (
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
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
		Description: "",
		ExeEnv: exeenv.ExeEnv{
			Env:     exeenv.InContainer,
			Check:   exeenv.InContainer,
			Exploit: exeenv.InContainer,
		},
		CheckSecPrerequisites:    nil,
		ExploitablePrerequisites: nil,
	},
}

func (v *vulnerability) Exploit(context *cli.Context) (err error) {
	if err := v.BaseVulnerability.Exploit(context); err != nil {
		return err
	}
	return Exploit()
}

func Exploit() (err error) {
	return Load()
}
