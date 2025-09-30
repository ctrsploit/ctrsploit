package apparmor

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_ProfileDockerDefault(t *testing.T) {
	tests := map[string]struct {
		satisfied bool
	}{
		"docker-v28.3.2": {
			satisfied: true,
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		satisfied, err := ProfileDockerDefault.Check()
		assert.NoError(t, err)
		assert.Equal(t, test.satisfied, satisfied)
	})
}
