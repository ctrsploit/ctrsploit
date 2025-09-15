package root

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/ssst0n3/awesome_libs/awesome_error"
)

func GetInodeNumber(path string) (ino uint64, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
			return
		}
		err = fmt.Errorf("unexpected error to stat %s: %w", path, err)
		awesome_error.CheckErr(err)
		return
	}
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		err = fmt.Errorf("failed to get syscall.Stat_t for %s", path)
		awesome_error.CheckErr(err)
		return
	}
	ino = stat.Ino
	return
}
