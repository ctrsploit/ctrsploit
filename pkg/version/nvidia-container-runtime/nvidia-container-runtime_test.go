package nvidia_container_runtime

import (
	"fmt"
	"os"
	"testing"

	"github.com/ctrsploit/ctrsploit/pkg/version/version"
	"github.com/stretchr/testify/assert"
)

func TestE2E_GetVersion(t *testing.T) {
	tests := map[string]struct {
		ver version.Number
	}{
		"nvidia-container-toolkit-v1.17.7": {
			ver: version.Number{
				Major: 1,
				Minor: 17,
				Patch: 7,
				Rc:    -1,
				Beta:  -1,
				Init:  true,
			},
		},
		"nvidia-container-toolkit-v1.17.0-rc.1": {
			ver: version.Number{
				Major: 1,
				Minor: 17,
				Patch: 0,
				Rc:    1,
				Beta:  -1,
				Init:  true,
			},
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
		assert.Equal(t, test.ver, *ver)
	})
}
