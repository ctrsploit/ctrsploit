package bpf

import (
	"github.com/ctrsploit/ctrsploit/vul/caps/sys_admin/ebpf"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

var (
	aliases = []string{"bpf"}
	VulCmd  = &cli.Command{
		Name:    Vul.GetName(),
		Aliases: aliases,
		Usage:   Vul.GetDescription(),
		Subcommands: []*cli.Command{
			ebpf.VulCmd,
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
			Description: "Container can load evil bpf program when has cap_bpf, may cause container escape",
			Level:       vul.LevelHigh,
			ExeEnv: exeenv.ExeEnv{
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
		},
	}
)
