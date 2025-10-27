package user

import (
	"os"

	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type EUidEqualTo struct {
	EUid int
	prerequisite.BasePrerequisite
}

func (p *EUidEqualTo) Check() (bool, error) {
	return p.CheckTemplate(func() (bool, error) {
		p.Satisfied = os.Geteuid() == p.EUid
		return p.Satisfied, p.Err
	})
}

var MustBeRootToWriteReleaseAgent = EUidEqualTo{
	EUid: 0,
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name: "euid=0",
		Info: "Current user must be root to write release_agent",
	},
}
