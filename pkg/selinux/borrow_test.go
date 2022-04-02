package selinux

import (
	"github.com/ssst0n3/awesome_libs/log"
	"testing"
)

func TestGetSelinuxMountPoint(t *testing.T) {
	mountPoint := GetSelinuxMountPoint()
	log.Logger.Info(mountPoint)
}
