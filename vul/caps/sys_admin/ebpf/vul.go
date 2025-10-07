package ebpf

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_admin/ebpf/bash"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_admin/ebpf/cron"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_admin/ebpf/execve"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

var (
	aliases     = []string{}
	CheckSecCmd = GetCheckSecCmd(Vul.GetName(), Vul.GetDescription(), aliases)
	ExploitCmd  = GetExploitCmd(Vul.GetName(), Vul.GetDescription(), aliases)
	VulCmd      = &cli.Command{
		Name:    Vul.GetName(),
		Aliases: aliases,
		Usage:   Vul.GetDescription(),
		Subcommands: []*cli.Command{
			GetCheckSecCmd("checksec", "check vulnerability exists", []string{"c"}),
			GetExploitCmd("exploit", "run the exploit", []string{"x"}),
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

var Vul = vulnerability{
	BaseVulnerability: vul.BaseVulnerability{
		Name:        "ebpf",
		Description: "",
		ExeEnv: exeenv.ExeEnv{
			Env:     exeenv.InContainer,
			Check:   exeenv.InContainer,
			Exploit: exeenv.InContainer,
		},
		CheckSecPrerequisites: prerequisite.Or(
			&capability.CapSysAdminBnd,
			prerequisite.And(
			// TODO: CapBpfAdmin, ...
			),
		),
		// TODO:
		ExploitablePrerequisites: nil,
	},
}

func GetCheckSecCmd(name, usage string, aliases []string) (cmd *cli.Command) {
	cmd = app.Vul2ChecksecCmd(&Vul, aliases, nil)
	cmd.Name = name
	cmd.Usage = usage
	cmd.Aliases = aliases
	return
}

func GetExploitCmd(name, usage string, aliases []string) *cli.Command {
	return &cli.Command{
		Name:    name,
		Aliases: aliases,
		Usage:   usage,
		Subcommands: []*cli.Command{
			execve.ExploitCmd,
			bash.ExploitCmd,
			cron.ExploitCmd,
		},
	}
}
