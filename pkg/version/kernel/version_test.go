package kernel

import (
	"fmt"
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

func TestE2E_Version(t *testing.T) {
	tests := map[string]struct {
		ver *semver.Version
	}{
		"cve-2022-0492": {
			ver: semver.New(5, 4, 0, "100-generic", ""),
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		ver, err := Version()
		assert.NoError(t, err)
		assert.Equal(t, test.ver, ver)
	})
}
