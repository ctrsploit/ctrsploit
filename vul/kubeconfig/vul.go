package kubeconfig

import (
	user_exec "github.com/ctrsploit/ctrsploit/vul/kubeconfig/user-exec"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v3"
)

var (
	aliases = []string{"kubecfg"}
	VulCmd  = &cli.Command{
		Name:    Vul.GetName(),
		Aliases: aliases,
		Usage:   Vul.GetDescription(),
		Commands: []*cli.Command{
			user_exec.VulCmd,
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

var Vul = vulnerability{
	BaseVulnerability: vul.BaseVulnerability{
		Name:        "kubeconfig",
		Description: "check kubeconfig related vulnerabilities",
		ExeEnv: exeenv.ExeEnv{
			Env:   exeenv.Local,
			Check: exeenv.Local,
		},
	},
}
