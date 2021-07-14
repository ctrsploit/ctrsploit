package util

import (
	"github.com/containerd/cgroups"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStat(t *testing.T) {
	control, err := cgroups.Load(cgroups.V1, cgroups.StaticPath(""))
	assert.NoError(t, err)
	metrics, err := control.Stat()
	assert.NoError(t, err)
	spew.Dump(metrics)
}
