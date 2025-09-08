package nvidia_container_runtime

import (
	"fmt"
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

func TestE2E_GetVersion(t *testing.T) {
	tests := map[string]struct {
		ver *semver.Version
	}{
		"nvidia-container-toolkit-v1.17.7": {
			ver: semver.New(1, 17, 7, "", ""),
		},
		"nvidia-container-toolkit-v1.17.0-rc.1": {
			ver: semver.New(1, 17, 0, "rc.1", ""),
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		ver, err := GetVersion()
		assert.NoError(t, err)
		assert.Equal(t, test.ver, ver)
	})
}
