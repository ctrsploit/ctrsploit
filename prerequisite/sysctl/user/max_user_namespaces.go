package user

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/sysctl"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type MaxUserNamespaces struct {
	prerequisite.BasePrerequisite
	GreaterThan int
}

func (p *MaxUserNamespaces) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	p.Checked = true

	maxUserNamespaces, err := sysctl.MaxUserNamespaces()
	if err != nil {
		return false, fmt.Errorf("unknown the meaning of: %w", err)
	}
	p.Satisfied = maxUserNamespaces > p.GreaterThan
	return p.Satisfied, nil
}

var (
	//goland:noinspection GoNameStartsWithPackageName
	UserNsEnabled = MaxUserNamespaces{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "user ns enabled",
			Info:   "user.max_user_namespaces > 0",
			ExeEnv: exeenv.InContainer | exeenv.InHost,
		},
		GreaterThan: 0,
	}
	//goland:noinspection GoNameStartsWithPackageName
	UserNsDisabled = prerequisite.Not(&UserNsEnabled)
)
