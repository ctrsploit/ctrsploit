package sys_admin

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/ctrsploit/vul/sys_admin/release_agent"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

var (
	aliases     = []string{"sys_admin"}
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

type SysAdmin struct {
	vul.BaseVulnerability
}

var (
	Vul = SysAdmin{
		BaseVulnerability: vul.BaseVulnerability{
			Name:        "cap_sys_admin",
			Description: "Container can be escaped when has cap_sys_admin",
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites:    &capability.CapSysAdminBnd,
			ExploitablePrerequisites: &capability.CapSysAdminEff,
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
			release_agent.ExploitCmd,
		},
	}
}
