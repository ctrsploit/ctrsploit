package pid

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/namespace"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_ptrace/pid_host"
	"github.com/ctrsploit/ctrsploit/vul/namespace/pid/proc_root"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v3"
)

var (
	aliases     = []string{"pid"}
	CheckSecCmd = getCheckSecCmd(Vul.Name, Vul.Description, aliases)
	ExploitCmd  = getExploitCmd(Vul.Name, Vul.Description, aliases)
	VulCmd      = &cli.Command{
		Name:    Vul.GetName(),
		Aliases: aliases,
		Usage:   Vul.GetDescription(),
		Commands: []*cli.Command{
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
			Name:        "host-pid",
			Description: "shared host pid namespace breaks process isolation",
			Level:       vul.LevelHigh,
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites:    &namespace.PidNamespaceLevelHost,
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
		Commands: []*cli.Command{
			proc_root.ExploitCmd,
			pid_host.GetExploitCmd(pid_host.Vul.GetName(), pid_host.Vul.GetDescription(), []string{"ptrace"}),
			// TODO: add more exploit methods
		},
	}
}
