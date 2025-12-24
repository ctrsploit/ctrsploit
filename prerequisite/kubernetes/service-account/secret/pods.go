package secret

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/kubernetes"
	kubernetesSecret "github.com/ctrsploit/ctrsploit/pkg/kubernetes/secret"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

// PodsWithSecretAccess is a prerequisite that checks for Pods with secret access permissions
// through their ServiceAccounts. It uses GetAllPodsWithSecretAccess to find all Pods
// that have access to secrets via RBAC bindings.
type PodsWithSecretAccess struct {
	prerequisite.BasePrerequisite
	pods []kubernetesSecret.PodWithSecretAccess
}

var HasPodsWithSecretAccess = PodsWithSecretAccess{
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name:   "HasPodsWithSecretAccess",
		Info:   "Check if there are any Pods with secret access permissions through ServiceAccounts",
		ExeEnv: exeenv.K8S,
	},
}

func (p *PodsWithSecretAccess) Check() (bool, error) {
	return p.CheckTemplate(func() {
		log.Logger.Debugf("Checking for Pods with secret access")
		clientset, err := kubernetes.GetKubernetesClient()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("failed to get Kubernetes client: %w", err))
			return
		}

		pods, err := kubernetesSecret.GetAllPodsWithSecretAccess(clientset)
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("failed to get pods with secret access: %w", err))
			return
		}

		p.pods = pods
		if len(pods) > 0 {
			p.Satisfied = true
			log.Logger.Debugf("Found %d pods with secret access", len(pods))
		} else {
			log.Logger.Debugf("No pods with secret access found")
		}
	})
}

// GetPodsWithSecretAccess extracts the list of Pods with secret access from a prerequisite.Set.
// It should be called after the prerequisite has been checked.
func GetPodsWithSecretAccess(s prerequisite.Set) ([]kubernetesSecret.PodWithSecretAccess, error) {
	_, _ = s.Check()
	if p, ok := s.(*PodsWithSecretAccess); ok {
		return p.pods, nil
	}
	if not, ok := s.(*prerequisite.SetNot); ok {
		if p, ok := not.Set.(*PodsWithSecretAccess); ok {
			return p.pods, nil
		}
		return nil, fmt.Errorf("unsupported prerequisite type inside Not: %T", not.Set)
	}
	return nil, fmt.Errorf("unsupported prerequisite type: %T", s)
}
