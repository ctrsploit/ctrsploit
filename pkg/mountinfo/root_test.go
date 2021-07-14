package mountinfo

import (
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRootMount(t *testing.T) {
	mount, err := RootMount()
	assert.NoError(t, err)
	log.Logger.Info(mount)
}
