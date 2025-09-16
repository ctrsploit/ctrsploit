package caps

import (
	"github.com/ctrsploit/ctrsploit/vul/caps/shocker"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_admin"
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_ptrace"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

var (
	aliases = []string{"caps"}
	VulCmd  = &cli.Command{
		Name:    Vul.GetName(),
		Aliases: aliases,
		Usage:   Vul.GetDescription(),
		Subcommands: []*cli.Command{
			sys_admin.VulCmd,
			sys_ptrace.VulCmd,
			shocker.VulCmd,
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

var (
	Vul = vulnerability{
		BaseVulnerability: vul.BaseVulnerability{
			Name:        "capability",
			Description: "Container can be dangerous when has some capabilities",
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
		},
	}
)
