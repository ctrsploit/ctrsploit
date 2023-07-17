package uname

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestUname(t *testing.T) {
	u, err := Uname()
	assert.NoError(t, err)
	spew.Dump(u)
}

func TestAll(t *testing.T) {
	cmd := exec.Command("uname", "-a")
	output, err := cmd.Output()
	assert.NoError(t, err)
	log.Logger.Info(string(output))
	log.Logger.Info(All())
}

func TestVersion(t *testing.T) {
	v, err := Version()
	assert.NoError(t, err)
	log.Logger.Info(v)
}

func TestRelease(t *testing.T) {
	r, err := Release()
	assert.NoError(t, err)
	log.Logger.Info(r)
}
