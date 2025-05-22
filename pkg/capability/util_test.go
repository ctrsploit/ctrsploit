package capability

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPid1Capability(t *testing.T) {
	caps, err := GetPid1Capability()
	assert.NoError(t, err)
	spew.Dump(caps)
}

func TestGetCurrentCapability(t *testing.T) {
	caps, err := GetCurrentCapability()
	assert.NoError(t, err)
	spew.Dump(caps)
}