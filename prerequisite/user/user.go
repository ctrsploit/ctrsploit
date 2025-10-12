package user

import (
	"os"

	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type EUidEqualTo struct {
	EUid int
	prerequisite.BasePrerequisite
}

func (p *EUidEqualTo) Check() (satisfied bool, err error) {
	if !p.Checked {
		p.Satisfied = os.Geteuid() == p.EUid
		p.Checked = true
	}
	satisfied = p.Satisfied
	return
}

var MustBeRootToWriteReleaseAgent = EUidEqualTo{
	EUid: 0,
	BasePrerequisite: prerequisite.BasePrerequisite{
		Name: "euid=0",
		Info: "Current user must be root to write release_agent",
	},
}
