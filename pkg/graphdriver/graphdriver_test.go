package graphdriver

import (
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGraphDriver_Init(t *testing.T) {
	g := GraphDriver{}
	assert.NoError(t, g.Init())
	log.Logger.Info(g.Type)
	log.Logger.Info(g.Rootfs)
}