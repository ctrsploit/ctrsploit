package hostname

import (
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"k8s.io/apimachinery/pkg/util/validation"
)

const K8sSafeAlphaNums = "bcdfghjklmnpqrstvwxz2456789"

type deployment struct {
	prerequisite.BasePrerequisite
	Pattern string
}

func (p *deployment) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	hostname, err := os.Hostname()
	if err != nil {
		return false, fmt.Errorf("failed to check %s, caused by unable to determine hostname: %w", p.Name, err)
	}
	p.Satisfied = p.check(hostname)
	p.Checked = true
	return p.Satisfied, nil
}

func (p *deployment) check(hostname string) bool {
	deploymentName, podTemplateHash, suffix, err := p.parse(hostname)
	if err != nil {
		// ignore error, maybe not deployment, maybe other container runtime
		return false
	}
	if !p.validSuffix(suffix) {
		return false
	}
	if !p.validHash(podTemplateHash) {
		return false
	}
	if !p.validDeploymentName(deploymentName, len(podTemplateHash)) {
		return false
	}
	return true
}

func (p *deployment) parse(hostname string) (deploymentName string, podTemplateHash string, suffix string, err error) {
	parts := strings.Split(hostname, "-")
	if len(parts) < 3 {
		err = fmt.Errorf("invalid hostname format: %q; expected format <deployment-name>-<pod-template-hash>-<suffix>", hostname)
		return
	}
	suffix = parts[len(parts)-1]
	podTemplateHash = parts[len(parts)-2]
	deploymentName = strings.Join(parts[:len(parts)-2], "-")
	return
}

// validSuffix
// https://github.com/kubernetes/kubernetes/blob/v1.34.1/staging/src/k8s.io/apiserver/pkg/storage/names/generate.go#L53
// https://github.com/kubernetes/kubernetes/blob/v1.34.1/staging/src/k8s.io/apimachinery/pkg/util/rand/rand.go#L83
func (p *deployment) validSuffix(suffix string) bool {
	if len(suffix) != 5 {
		return false
	}
	for _, r := range suffix {
		if !strings.ContainsRune(K8sSafeAlphaNums, r) {
			return false
		}
	}
	return true
}

// validHash
// https://github.com/kubernetes/kubernetes/blob/v1.34.1/pkg/controller/controller_utils.go#L1382
// https://github.com/kubernetes/kubernetes/blob/v1.34.1/pkg/controller/deployment/sync.go#L561
func (p *deployment) validHash(hash string) bool {
	// sum32
	if len(hash) < 1 || len(hash) > 10 {
		return false
	}
	for _, r := range hash {
		if !strings.ContainsRune(K8sSafeAlphaNums, r) {
			return false
		}
	}
	return true
}

// validDeploymentName
// a DNS-1123 subdomain must consist of lower case alphanumeric characters, '-' or '.',
// and must start and end with an alphanumeric character (e.g. 'example.com',
// regex used for validation is '[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*')
func (p *deployment) validDeploymentName(name string, hashLen int) bool {
	replicaSetNameSeparator := "-"
	maxDeploymentNameLength := validation.DNS1123SubdomainMaxLength - len(replicaSetNameSeparator) - hashLen
	if len(name) > maxDeploymentNameLength {
		return false
	}
	// https://github.com/kubernetes/kubernetes/blob/master/pkg/apis/apps/validation/validation.go#L543
	errs := validation.IsDNS1123Subdomain(name)
	if len(errs) > 0 {
		return false
	}
	return true
}

var (
	K8sDeploymentHostname = deployment{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "hostname",
			Info:   "k8s hostname match RFC 1123",
			ExeEnv: exeenv.InContainer,
		},
		Pattern: "^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*-[0-9a-f]{9,10}-[bcdfghjklmnpqrstvwxz2456789]{5}$",
	}
)
