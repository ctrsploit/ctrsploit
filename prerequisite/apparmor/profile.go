package apparmor

import (
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type ProfileNameContains struct {
	prerequisite.BasePrerequisite
	Expected string
}

func (p *ProfileNameContains) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	profile, err := os.ReadFile("/proc/self/attr/apparmor/current")
	if err != nil {
		return false, fmt.Errorf("could not read apparmor profile: %w", err)
	}
	p.Satisfied = strings.Contains(string(profile), p.Expected)
	p.Checked = true
	return p.Satisfied, nil
}

func NewProfileNameContains(name string) *ProfileNameContains {
	return &ProfileNameContains{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   fmt.Sprintf("apparmor %s", name),
			Info:   fmt.Sprintf("apparmor profile is %s", name),
			ExeEnv: exeenv.InContainer,
		},
		Expected: name,
	}
}

var (
	// ProfileDockerDefault https://github.com/moby/moby/blob/v28.4.0/daemon/apparmor_default.go#L15
	ProfileDockerDefault = NewProfileNameContains("docker-default")
	// ProfileNerdctlDefault https://github.com/containerd/nerdctl/blob/v2.1.5/pkg/defaults/defaults_linux.go#L32
	ProfileNerdctlDefault = NewProfileNameContains("nerdctl-default")
	// ProfileCriContainerd https://github.com/containerd/containerd/blob/v2.1.4/internal/cri/sputil/securityprofile_linux.go#L34
	ProfileCriContainerd = NewProfileNameContains("cri-containerd.apparmor.d")
)
