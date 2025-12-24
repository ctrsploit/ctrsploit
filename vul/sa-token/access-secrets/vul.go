package access_secrets

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/prerequisite/kubernetes/service-account/secret"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

var (
	aliases = []string{"secret"}

	checkSecFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    "kubeconfig",
			Aliases: []string{"k"},
			Usage:   "Path to kubeconfig file (defaults to in-cluster config, then ~/.kube/config)",
			EnvVars: []string{"KUBECONFIG"},
		},
	}

	CheckSecCmd = getCheckSecCmd(Vul.Name, Vul.Description, aliases)
	VulCmd      = &cli.Command{
		Name:    Vul.Name,
		Aliases: aliases,
		Usage:   Vul.Description,
		Subcommands: []*cli.Command{
			getCheckSecCmd("checksec", "check vulnerability exists", []string{"c"}),
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

var Vul = vulnerability{
	BaseVulnerability: vul.BaseVulnerability{
		Name:        "sa-token-access-secrets",
		Description: "Check if service account token can access Kubernetes Secrets",
		ExeEnv: exeenv.ExeEnv{
			Env:   exeenv.K8S,
			Check: exeenv.K8S,
		},
		CheckSecPrerequisites:    &secret.HasPodsWithSecretAccess,
		ExploitablePrerequisites: nil,
	},
}

func (v *vulnerability) CheckSec(ctx *cli.Context) (satisfied bool, err error) {
	log.Logger.Debugf("Starting vulnerability.CheckSec for service account token secrets access")

	// Check prerequisites first
	satisfied, err = v.BaseVulnerability.CheckSec(ctx)
	if err != nil {
		return false, fmt.Errorf("prerequisite check failed: %w", err)
	}
	if !satisfied {
		log.Logger.Info("Prerequisites not satisfied: no pods with secret access found")
		return false, nil
	}

	// Report the results using the exported function
	if err := Check(); err != nil {
		return false, err
	}

	return true, nil
}

func getCheckSecCmd(name, usage string, aliases []string) (cmd *cli.Command) {
	cmd = app.Vul2ChecksecCmd(&Vul, aliases, checkSecFlags)
	cmd.Name = name
	cmd.Usage = usage
	return
}
