package group

import (
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os/user"
	"strconv"
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

func (p *MustBe) Check() (err error) {
	err = p.BasePrerequisite.Check()
	if err != nil {
		return
	}
	current, err := user.Current()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	gid, err := strconv.Atoi(current.Gid)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	p.Satisfied = uint(gid) == p.ExpectedGroup
	return
}
