package seccomp

import (
	"github.com/ssst0n3/awesome_libs/log"
	"testing"
)

func TestCheckSupported(t *testing.T) {
	log.Logger.Info(CheckSupported())
}

func TestCheckEnabled(t *testing.T) {
	log.Logger.Info(CheckEnabled())
}