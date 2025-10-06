package rootfs

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_DockerOverlay(t *testing.T) {
	tests := map[string]struct {
		is          bool
		pathPattern string
	}{
		"docker-v28.3.2": {
			is:          true,
			pathPattern: "^/var/lib/docker/overlay2/[a-f0-9]{64}/merged$",
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		d := DockerOverlay{}
		is, err := d.Is()
		assert.NoError(t, err)
		assert.Equal(t, test.is, is)
		path, err := d.RootPath()
		assert.NoError(t, err)
		r, err := regexp.Compile(test.pathPattern)
		assert.NoError(t, err)
		assert.True(t, r.MatchString(path))
	})
}
