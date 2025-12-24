package access_secrets

import (
	"fmt"
	"strings"

	sa_secret "github.com/ctrsploit/ctrsploit/prerequisite/kubernetes/service-account/secret"
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
		CheckSecPrerequisites:    &sa_secret.HasPodsWithSecretAccess,
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

	// Get pods with secret access from the prerequisite
	pods, err := sa_secret.GetPodsWithSecretAccess(&sa_secret.HasPodsWithSecretAccess)
	if err != nil {
		return false, fmt.Errorf("failed to get pods with secret access: %w", err)
	}

	if len(pods) == 0 {
		log.Logger.Info("No pods with secret access found")
		return false, nil
	}

	// Collect cluster-wide statistics
	allNamespaces := make(map[string]bool)
	allPermissions := make(map[string]bool)
	namespacePodCount := make(map[string]int)
	serviceAccountMap := make(map[string]map[string]bool) // namespace -> serviceAccount -> true

	for _, pod := range pods {
		allNamespaces[pod.Namespace] = true
		namespacePodCount[pod.Namespace]++

		// Track service accounts
		if serviceAccountMap[pod.Namespace] == nil {
			serviceAccountMap[pod.Namespace] = make(map[string]bool)
		}
		serviceAccountMap[pod.Namespace][pod.ServiceAccountName] = true

		// Track unique permissions
		for _, perm := range pod.Permissions {
			permKey := fmt.Sprintf("%s/%s/%s", perm.BindingKind, perm.RoleRefKind, perm.RoleRefName)
			allPermissions[permKey] = true
		}
	}

	// Display cluster-wide information
	log.Logger.Infof("")
	log.Logger.Infof("=== Cluster-wide Pods with Secret Access ===")
	log.Logger.Infof("Found %d pod(s) with secret access permissions across the cluster", len(pods))

	// Display pods grouped by namespace
	log.Logger.Infof("")
	log.Logger.Infof("=== Pods by Namespace ===")
	for namespace := range allNamespaces {
		log.Logger.Infof("")
		log.Logger.Infof("Namespace: %s", namespace)
		log.Logger.Infof("  Pods: %d", namespacePodCount[namespace])
		log.Logger.Infof("  Service Accounts: %d", len(serviceAccountMap[namespace]))

		// Show unique RBAC permissions for this namespace
		seenPerms := make(map[string]bool)
		for _, pod := range pods {
			if pod.Namespace == namespace {
				for _, perm := range pod.Permissions {
					permKey := fmt.Sprintf("%s/%s/%s/%s", perm.BindingKind, perm.RoleRefKind, perm.RoleRefName, strings.Join(perm.Verbs, ","))
					if !seenPerms[permKey] {
						seenPerms[permKey] = true
						if len(seenPerms) == 1 {
							log.Logger.Infof("  RBAC Permissions:")
						}
						log.Logger.Infof("    - [%s] %s (Scope: %s, Verbs: %s)", perm.BindingKind, perm.RoleRefName, perm.Scope(), strings.Join(perm.Verbs, ","))
					}
				}
			}
		}
	}

	// Display summary
	log.Logger.Infof("")
	log.Logger.Infof("=== Summary ===")
	log.Logger.Infof("Total pods with secret access: %d", len(pods))
	log.Logger.Infof("Namespaces affected: %d", len(allNamespaces))
	log.Logger.Infof("Unique RBAC permissions: %d", len(allPermissions))

	return true, nil
}

func getCheckSecCmd(name, usage string, aliases []string) (cmd *cli.Command) {
	cmd = app.Vul2ChecksecCmd(&Vul, aliases, checkSecFlags)
	cmd.Name = name
	cmd.Usage = usage
	return
}
