package where

import (
	"ctrsploit/test/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckDockerEnvExists(t *testing.T) {
	d := Docker{}
	d.CheckDockerEnvExists()
	if config.InDocker {
		assert.Equal(t, true, d.DockerEnvFileExists)
	} else {
		assert.Equal(t, false, d.DockerEnvFileExists)
	}
}

func TestCheckMountInfo(t *testing.T) {
	d := Docker{}
	assert.NoError(t, d.CheckMountInfo())
	if config.InDocker {
		assert.Equal(t, true, d.RootfsContainsDocker)
	} else {
		assert.Equal(t, false, d.RootfsContainsDocker)
	}
}
