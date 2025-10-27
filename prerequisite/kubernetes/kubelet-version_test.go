package kubernetes

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestE2E_FixedCVE_2020_8558(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	tests := map[string]struct {
		expected bool
	}{
		"kubernetes-v1.34.0-host": {
			true,
		},
		"kubernetes-v1.18.2-host": {
			false,
		},
	}
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(testEnv, func(t *testing.T) {
		satisfied, err := FixedCVE_2020_8558.Check()
		require.NoError(t, err)
		require.Equal(t, test.expected, satisfied)
	})
}
