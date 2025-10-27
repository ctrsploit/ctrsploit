package nvidia_container_toolkit

import (
	"github.com/Masterminds/semver/v3"
	nvidia_container_runtime "github.com/ctrsploit/ctrsploit/pkg/version/nvidia-container-runtime"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type Version struct {
	prerequisite.BasePrerequisite
	Constraint *semver.Constraints
}

func (p *Version) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	ver, err := nvidia_container_runtime.GetVersion()
	if err != nil {
		return false, err
	}
	p.Satisfied = p.Constraint.Check(ver)
	p.Checked = true
	return p.Satisfied, nil
}

var (
	VulnerableTo0132 = func() *Version {
		cons, err := semver.NewConstraint(">= 1.0.0, <= 1.16.1")
		if err != nil {
			panic(err)
		}
		return &Version{
			BasePrerequisite: prerequisite.BasePrerequisite{
				Name:   "nvidia-container-toolkit vulnerable to CVE-2024-0132",
				Info:   "nvidia-container-toolkit(libnvidia-container) >=v1.0.0, <=v1.16.1 is vulnerable to CVE-2024-0132",
				ExeEnv: exeenv.InHost | exeenv.Local,
			},
			Constraint: cons,
		}
	}()
	VulnerableTo23266 = func() *Version {
		cons, err := semver.NewConstraint(">= 1.10.0-rc.1, <= 1.17.7")
		if err != nil {
			panic(err)
		}
		return &Version{
			BasePrerequisite: prerequisite.BasePrerequisite{
				Name:   "nvidia-container-toolkit vulnerable to CVE-2025-23266",
				Info:   "nvidia-container-toolkit >= 1.10.0-rc.1, <= 1.17.7 is vulnerable to CVE-2025-23266",
				ExeEnv: exeenv.InHost | exeenv.Local,
			},
			Constraint: cons,
		}
	}()
)
