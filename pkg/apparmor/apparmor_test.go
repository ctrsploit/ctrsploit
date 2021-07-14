package apparmor

import (
	"github.com/ssst0n3/awesome_libs/log"
	"testing"
)

func TestIsSupport(t *testing.T) {
	log.Logger.Info(IsSupport())
}
