package group

import (
	"fmt"
	"os/user"
	"strconv"

	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type MustBe struct {
	ExpectedGroup uint
	prerequisite.BasePrerequisite
}

var MustBeRoot = MustBe{
	ExpectedGroup: 0,
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name: "root",
		Info: "Current group must be root",
	},
}

func (p *MustBe) Check() (bool, error) {
	return p.CheckTemplate(func() (bool, error) {
		current, err := user.Current()
		if err != nil {
			p.Err = fmt.Errorf("failed to check [%s], caused by getting current user: %w", p.GetName(), err)
			return p.Satisfied, p.Err
		}
		gid, err := strconv.Atoi(current.Gid)
		if err != nil {
			p.Err = fmt.Errorf("failed to check [%s], caused by converting gid to int: %w", p.GetName(), err)
			return p.Satisfied, p.Err
		}
		p.Satisfied = uint(gid) == p.ExpectedGroup
		return p.Satisfied, p.Err
	})
}
