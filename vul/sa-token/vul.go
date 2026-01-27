package caps

import (
	access_secrets "github.com/ctrsploit/ctrsploit/vul/sa-token/access-secrets"
	"github.com/ctrsploit/ctrsploit/vul/sa-token/policy"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v3"
)

var (
	aliases = []string{"sa-token", "token"}
	VulCmd  = &cli.Command{
		Name:    Vul.GetName(),
		Aliases: aliases,
		Usage:   Vul.GetDescription(),
		Commands: []*cli.Command{
			access_secrets.VulCmd,
			policy.VulCmd,
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

var (
	Vul = vulnerability{
		BaseVulnerability: vul.BaseVulnerability{
			Name:        "service-account-token",
			Description: "check service account token related vulnerabilities",
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.K8S,
				Check:   exeenv.K8S,
				Exploit: exeenv.K8S,
			},
		},
	}
)
