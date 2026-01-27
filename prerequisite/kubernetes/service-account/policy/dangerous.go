package policy

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/kubernetes"
	kubernetesPolicy "github.com/ctrsploit/ctrsploit/pkg/kubernetes/policy"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

// DangerousPermissions is a prerequisite that checks if current ServiceAccount token
// has dangerous RBAC permissions that could lead to privilege escalation or RCE.
type DangerousPermissions struct {
	prerequisite.BasePrerequisite
	results     []kubernetesPolicy.CheckResult
	permissions []kubernetesPolicy.DangerousPermission
	namespace   string
	minLevel    kubernetesPolicy.Level
}

var HasDangerousPermissions = DangerousPermissions{
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name:   "HasDangerousPermissions",
		Info:   "Check if ServiceAccount token has dangerous RBAC permissions",
		ExeEnv: exeenv.InContainer | exeenv.K8S,
	},
	permissions: kubernetesPolicy.DefaultPermissions,
	namespace:   "", // empty means cluster-wide
	minLevel:    kubernetesPolicy.LevelMedium,
}

// NewDangerousPermissions creates a new DangerousPermissions prerequisite with custom settings
func NewDangerousPermissions(namespace string, minLevel kubernetesPolicy.Level, permissions []kubernetesPolicy.DangerousPermission) *DangerousPermissions {
	if permissions == nil {
		permissions = kubernetesPolicy.DefaultPermissions
	}
	return &DangerousPermissions{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "HasDangerousPermissions",
			Info:   "Check if ServiceAccount token has dangerous RBAC permissions",
			ExeEnv: exeenv.InContainer | exeenv.K8S,
		},
		permissions: permissions,
		namespace:   namespace,
		minLevel:    minLevel,
	}
}

func (p *DangerousPermissions) Check() (bool, error) {
	return p.CheckTemplate(func() {
		log.Logger.Debugf("Checking for dangerous permissions")
		clientset, err := kubernetes.GetKubernetesClient()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("failed to get Kubernetes client: %w", err))
			return
		}

		// Filter permissions by minimum level
		permissions := kubernetesPolicy.FilterByLevel(p.permissions, p.minLevel)

		results, err := kubernetesPolicy.CheckDangerousPermissions(clientset, p.namespace, permissions)
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("failed to check dangerous permissions: %w", err))
			return
		}

		p.results = results
		if len(results) > 0 {
			p.Satisfied = true
			log.Logger.Debugf("Found %d dangerous permissions", len(results))
		} else {
			log.Logger.Debugf("No dangerous permissions found")
		}
	})
}

// GetResults returns the check results after Check() has been called
func (p *DangerousPermissions) GetResults() []kubernetesPolicy.CheckResult {
	return p.results
}

// GetDangerousPermissionResults extracts the check results from a prerequisite.Set.
// It should be called after the prerequisite has been checked.
func GetDangerousPermissionResults(s prerequisite.Set) ([]kubernetesPolicy.CheckResult, error) {
	_, _ = s.Check()
	if p, ok := s.(*DangerousPermissions); ok {
		return p.results, nil
	}
	if not, ok := s.(*prerequisite.SetNot); ok {
		if p, ok := not.Set.(*DangerousPermissions); ok {
			return p.results, nil
		}
		return nil, fmt.Errorf("unsupported prerequisite type inside Not: %T", not.Set)
	}
	return nil, fmt.Errorf("unsupported prerequisite type: %T", s)
}
