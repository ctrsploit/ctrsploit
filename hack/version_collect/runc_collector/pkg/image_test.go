package pkg

import (
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_inspect(t *testing.T) {
	inspectGraphDrivers("docker:24.0.7-dind-alpine3.19")
}

// TestGetFileHostPath
// test: sudo -E /home/st0n3/sdk/go1.21.5/bin/go test -test.run ^TestGetFileHostPath$
func TestGetFileHostPath(t *testing.T) {
	GetFileHostPath("docker:24.0.7-dind-alpine3.19", "/usr/local/bin/runc")
}

func Test_ListTags(t *testing.T) {
	tags, err := ListTags("docker")
	assert.NoError(t, err)
	log.Logger.Info(tags)
}

func TestPull(t *testing.T) {
	Pull("docker:1")
}
