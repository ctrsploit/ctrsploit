package noused

import (
	"github.com/ctrsploit/ctrsploit/test/config"
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetContainerFullId(t *testing.T) {
	if !config.InDocker {
		t.Skip("container id lookup requires Docker cgroup metadata")
	}
	id, err := GetContainerFullId()
	assert.NoError(t, err)
	log.Logger.Debug(id)
	assert.Equal(t, true, len(id) > 0)
}
