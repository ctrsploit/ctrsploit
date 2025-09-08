package group

import (
	"os/user"
	"strconv"

	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/awesome_error"
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
	if p.Checked {
		return p.Satisfied, nil
	}
	current, err := user.Current()
	if err != nil {
		awesome_error.CheckErr(err)
		return false, err
	}
	gid, err := strconv.Atoi(current.Gid)
	if err != nil {
		awesome_error.CheckErr(err)
		return false, err
	}
	p.Satisfied = uint(gid) == p.ExpectedGroup
	p.Checked = true
	return p.Satisfied, nil
}
