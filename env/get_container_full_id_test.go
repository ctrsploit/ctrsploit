package env

import (
	"ctrsploit/test/config"
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetContainerFullId(t *testing.T) {
	id, err := GetContainerFullId()
	assert.NoError(t, err)
	log.Logger.Debug(id)
	if config.InDocker {
		assert.Equal(t, true, len(id) > 0)
	}
}
