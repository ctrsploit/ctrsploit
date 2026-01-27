package pids

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_UnlimitedPidsMax(t *testing.T) {
	tests := map[string]struct {
		satisfied bool
	}{
		"unlimited": {
			satisfied: true,
		},
		"limited": {
			satisfied: false,
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		satisfied, err := UnlimitedPidsMax.Check()
		require.NoError(t, err)
		assert.Equal(t, test.satisfied, satisfied)
	})
}
