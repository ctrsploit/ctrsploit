package apparmor

import (
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type ProfileNameEqualTo struct {
	prerequisite.BasePrerequisite
	Name string
}

func (p *ProfileNameEqualTo) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	profile, err := os.ReadFile("/proc/self/attr/apparmor/current")
	if err != nil {
		return false, fmt.Errorf("could not read apparmor profile: %w", err)
	}
	p.Satisfied = strings.Contains(string(profile), p.Name)
	p.Checked = true
	return p.Satisfied, nil
}

func NewProfileNameEqualTo(name string) *ProfileNameEqualTo {
	return &ProfileNameEqualTo{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   fmt.Sprintf("apparmor %s", name),
			Info:   fmt.Sprintf("apparmor profile is %s", name),
			ExeEnv: exeenv.InContainer,
		},
		Name: name,
	}
}

var (
	// ProfileDockerDefault https://github.com/moby/moby/blob/v28.4.0/daemon/apparmor_default.go#L15
	ProfileDockerDefault = NewProfileNameEqualTo("docker-default")
	// ProfileNerdctlDefault https://github.com/containerd/nerdctl/blob/v2.1.5/pkg/defaults/defaults_linux.go#L32
	ProfileNerdctlDefault = NewProfileNameEqualTo("nerdctl-default")
)
