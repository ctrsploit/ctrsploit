package release_agent

import (
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/urfave/cli/v2"
)

var (
	aliases      = []string{"ra", "cgroup_v1_release_agent"}
	flagsExploit = []cli.Flag{
		&cli.StringFlag{Name: "cmd", Aliases: []string{"c"}, Required: true},
	}
	CheckSecCmd = app.Vul2ChecksecCmd(&Vul, aliases, nil)
	ExploitCmd  = app.Vul2ExploitCmd(&Vul, aliases, flagsExploit, true)
	VulCmd      = app.Vul2VulCmd(&Vul, aliases, nil, flagsExploit, true)
)
