package bpf

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_admin/ebpf"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

var (
	aliases     = []string{"bpf"}
	CheckSecCmd = getCheckSecCmd(Vul.GetName(), Vul.GetDescription(), aliases)
	ExploitCmd  = getExploitCmd(Vul.GetName(), Vul.GetDescription(), aliases)
	VulCmd      = &cli.Command{
		Name:    Vul.GetName(),
		Aliases: aliases,
		Usage:   Vul.GetDescription(),
		Subcommands: []*cli.Command{
			getCheckSecCmd("checksec", "check vulnerability exists", []string{"c"}),
			getExploitCmd("exploit", "run the exploit", []string{"x"}),
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

var (
	Vul = vulnerability{
		BaseVulnerability: vul.BaseVulnerability{
			Name:        "cap_bpf",
			Description: "load evil bpf programs via cap_bpf",
			Level:       vul.LevelHigh,
			ExeEnv: exeenv.ExeEnv{
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites: prerequisite.Or(
				&capability.CapBpfBnd,
			),
			ExploitablePrerequisites: nil,
		},
	}
)

func getCheckSecCmd(name, usage string, aliases []string) (cmd *cli.Command) {
	cmd = app.Vul2ChecksecCmd(&Vul, aliases, nil)
	cmd.Name = name
	cmd.Usage = usage
	cmd.Aliases = aliases
	return
}

func getExploitCmd(name, usage string, aliases []string) (cmd *cli.Command) {
	return &cli.Command{
		Name:    name,
		Usage:   usage,
		Aliases: aliases,
		Subcommands: []*cli.Command{
			// TODO: add more exploit methods
			ebpf.GetExploitCmd("ebpf", "exploit via eBPF", nil),
		},
	}
}
