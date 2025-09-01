package buildkitd

import (
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

func TestE2E_Version(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	tests := map[string]struct {
		name     string
		expected *semver.Version
	}{
		"buildkit-v0.12.4": {
			name:     "v0.12.4",
			expected: semver.MustParse("v0.12.4"),
		},
		"buildkit-v0.9.0": {
			name:     "v0.9.0",
			expected: semver.MustParse("v0.9.0"),
		},
	}
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(test.name, func(t *testing.T) {
		version, err := Version("unix:///var/run/buildkit/buildkitd.sock")
		assert.NoError(t, err)
		assert.Equal(t, test.expected, version)
	})
}
