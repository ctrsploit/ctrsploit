package policy

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/kubernetes/policy"
	prerequisitePolicy "github.com/ctrsploit/ctrsploit/prerequisite/kubernetes/service-account/policy"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v3"
)

var (
	aliases = []string{"policy", "dangerous-permissions", "dp"}

	checkSecFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    "namespace",
			Aliases: []string{"n"},
			Usage:   "Namespace to check (default: cluster-wide)",
		},
		&cli.StringFlag{
			Name:    "level",
			Aliases: []string{"l"},
			Usage:   "Minimum level to report: critical, high, medium (default: medium)",
			Value:   "medium",
		},
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to custom dangerous permissions YAML file (optional)",
		},
	}

	CheckSecCmd = getCheckSecCmd(Vul.Name, Vul.Description, aliases)
	VulCmd      = &cli.Command{
		Name:    Vul.Name,
		Aliases: aliases,
		Usage:   Vul.Description,
		Commands: []*cli.Command{
			getCheckSecCmd("checksec", "check vulnerability exists", []string{"c"}),
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

var Vul = vulnerability{
	BaseVulnerability: vul.BaseVulnerability{
		Name:        "sa-token-policy",
		Description: "Check if service account token has dangerous permissions",
		ExeEnv: exeenv.ExeEnv{
			Env:   exeenv.InContainer | exeenv.K8S,
			Check: exeenv.InContainer | exeenv.K8S,
		},
		CheckSecPrerequisites:    &prerequisitePolicy.HasDangerousPermissions,
		ExploitablePrerequisites: nil,
	},
}

func (v *vulnerability) CheckSec(cmd *cli.Command) (satisfied bool, err error) {
	log.Logger.Debugf("Starting vulnerability.CheckSec for service account token dangerous permissions")

	// Get parameters from flags
	namespace := cmd.String("namespace")
	levelStr := cmd.String("level")
	configPath := cmd.String("config")

	// Parse level
	minLevel := policy.LevelMedium
	switch levelStr {
	case "critical":
		minLevel = policy.LevelCritical
	case "high":
		minLevel = policy.LevelHigh
	case "medium":
		minLevel = policy.LevelMedium
	}

	// Load permissions
	permissions := policy.DefaultPermissions
	if configPath != "" {
		customPerms, err := policy.LoadPermissionsFromFile(configPath)
		if err != nil {
			return false, fmt.Errorf("failed to load custom permissions: %w", err)
		}
		permissions = policy.MergePermissions(permissions, customPerms)
	}

	// Create custom prerequisite with parameters
	prereq := prerequisitePolicy.NewDangerousPermissions(namespace, minLevel, permissions)

	// Check prerequisites
	satisfied, err = prereq.Check()
	if err != nil {
		return false, fmt.Errorf("prerequisite check failed: %w", err)
	}
	if !satisfied {
		log.Logger.Info("No dangerous permissions found")
		return false, nil
	}

	// Report the results
	results := prereq.GetResults()
	if err := Check(results); err != nil {
		return false, err
	}

	// Set the VulnerabilityExists field so that Output() displays [Y]
	v.VulnerabilityExists = true

	return true, nil
}

func getCheckSecCmd(name, usage string, cmdAliases []string) (cmd *cli.Command) {
	cmd = app.Vul2ChecksecCmd(&Vul, cmdAliases, checkSecFlags)
	cmd.Name = name
	cmd.Usage = usage
	return
}
