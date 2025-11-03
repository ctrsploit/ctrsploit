package env

import (
	"os"

	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type Exists struct {
	prerequisite.BasePrerequisite
	Expected string
}

func (p *Exists) Check() (bool, error) {
	return p.CheckTemplate(func() {
		_, p.Satisfied = os.LookupEnv(p.Expected)
		return
	})
}

var (
	// KubernetesServiceHostExists
	// https://github.com/kubernetes/kubernetes/blob/v1.34.1/pkg/kubelet/envvars/envvars.go#L45
	KubernetesServiceHostExists = Exists{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "KUBERNETES_SERVICE_HOST",
			Info:   "env KUBERNETES_SERVICE_HOST exists",
			ExeEnv: exeenv.InContainer,
		},
		Expected: "KUBERNETES_SERVICE_HOST",
	}
)
