package selinux

import (
	"github.com/ssst0n3/awesome_libs/log"
	"testing"
)

func TestIsEnabled(t *testing.T) {
	log.Logger.Info(IsEnabled())
}
