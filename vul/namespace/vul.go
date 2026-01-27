package namespace

import (
	"github.com/ctrsploit/ctrsploit/vul/namespace/net"
	"github.com/ctrsploit/ctrsploit/vul/namespace/pid"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v3"
)

var (
	aliases = []string{"ns"}
	VulCmd  = &cli.Command{
		Name:    Vul.GetName(),
		Aliases: aliases,
		Usage:   Vul.GetDescription(),
		Commands: []*cli.Command{
			pid.VulCmd,
			net.VulCmd,
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

var (
	Vul = vulnerability{
		BaseVulnerability: vul.BaseVulnerability{
			Name:        "namespace",
			Description: "host level namespaces break the isolations",
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
		},
	}
)
