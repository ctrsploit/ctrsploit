package user

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type MaxUserNamespaces struct {
	prerequisite.BasePrerequisite
	ExpectZero bool
	// TODO: add a related member
}

func (p *MaxUserNamespaces) Check() (bool, error) {
	if p.Checked {
		return p.Satisfied, nil
	}
	p.Checked = true
	content, err := os.ReadFile("/proc/sys/user/max_user_namespaces")
	if err != nil {
		if os.IsNotExist(err) {
			p.Satisfied = p.ExpectZero == true
			return p.Satisfied, nil
		} else {
			return false, fmt.Errorf("unknown the meaning of: %w", err)
		}
	}
	value, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		return false, fmt.Errorf("[NOT POSSIBLE] failed to parse the value of max_user_namespaces: %w", err)
	}
	p.Satisfied = p.ExpectZero == (value == 0)
	return p.Satisfied, nil
}

var (
	UserNsEnabled = MaxUserNamespaces{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "user ns enabled",
			Info:   "/proc/sys/user/max_user_namespaces > 0",
			ExeEnv: exeenv.InContainer | exeenv.InHost,
		},
		ExpectZero: false,
	}
	UserNsDisabled = MaxUserNamespaces{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   "user ns disabled",
			Info:   "/proc/sys/user/max_user_namespaces = 0",
			ExeEnv: exeenv.InContainer | exeenv.InHost,
		},
		ExpectZero: true,
	}
)
