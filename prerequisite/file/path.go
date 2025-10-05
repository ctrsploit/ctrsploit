package file

import (
	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type Exists struct {
	prerequisite.BasePrerequisite
	Path string
}

func (p *Exists) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	p.Satisfied, _ = internal.CheckPathExists(p.Path)
	p.Checked = true
	return p.Satisfied, nil
}

var (
	DockerEnvFileExists = Exists{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   ".dockerenv",
			Info:   "/.dockerenv exists",
			ExeEnv: exeenv.InContainer,
		},
		Path: "/.dockerenv",
	}
	// K8sSecretsExists
	// https://github.com/kubernetes/kubernetes/blob/v1.34.1/plugin/pkg/admission/serviceaccount/admission.go#L57
	K8sSecretsExists = Exists{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "secrets",
			Info:   "/var/run/secrets/kubernetes.io exists",
			ExeEnv: exeenv.InContainer,
		},
		Path: "/var/run/secrets/kubernetes.io",
	}
)
