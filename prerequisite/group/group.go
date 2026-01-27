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
	return p.CheckTemplate(func() {
		current, err := user.Current()
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("getting current user: %w", err))
			return
		}
		gid, err := strconv.Atoi(current.Gid)
		if err != nil {
			p.Err = p.WrapErr(fmt.Errorf("converting gid to int: %w", err))
			return
		}
		p.Satisfied = uint(gid) == p.ExpectedGroup
		return
	})
}
