package prerequisite

import (
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os/user"
	"strconv"
)

type UserMustBe struct {
	ExpectedUser uint
	BasePrerequisite
}

var MustBeRoot = UserMustBe{
	ExpectedUser: 0,
	BasePrerequisite: BasePrerequisite{
		Name: "root",
		Info: "Current user must be root",
	},
}

var MustBeRootToWriteReleaseAgent = UserMustBe{
	ExpectedUser: MustBeRoot.ExpectedUser,
	BasePrerequisite: BasePrerequisite{
		Name: MustBeRoot.Name,
		Info: "Current user must be root to write release_agent",
	},
}

func (p UserMustBe) Check() (err error) {
	current, err := user.Current()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	u, err := strconv.Atoi(current.Uid)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	p.Satisfied = uint(u) == p.ExpectedUser
	return
}
