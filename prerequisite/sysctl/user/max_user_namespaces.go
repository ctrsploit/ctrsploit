package user

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/sysctl"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type MaxUserNamespaces struct {
	prerequisite.BasePrerequisite
	GreaterThan uint64
}

func (p *MaxUserNamespaces) Check() (bool, error) {
	return p.CheckTemplate(func() (bool, error) {
		maxUserNamespaces, err := sysctl.MaxUserNamespaces()
		if err != nil {
			p.Err = fmt.Errorf("failed to check [%s], caused by getting sysctl: %w", p.GetName(), err)
			return p.Satisfied, p.Err
		}
		p.Satisfied = maxUserNamespaces > p.GreaterThan
		return p.Satisfied, p.Err
	})
}

//goland:noinspection GoNameStartsWithPackageName
var (
	UserNsEnabled = MaxUserNamespaces{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "user ns enabled",
			Info:   "user.max_user_namespaces > 0",
			ExeEnv: exeenv.InContainer | exeenv.InHost,
		},
		GreaterThan: 0,
	}
	UserNsDisabled = prerequisite.Not(&UserNsEnabled)
)
