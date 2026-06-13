package module

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestE2E_ModuleLoaded(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	tests := map[string][]struct {
		name      string
		module    string
		satisfied bool
	}{
		"docker-v28.3.2": {
			{
				name:      "loaded overlay module",
				module:    "overlay",
				satisfied: true,
			},
			{
				name:      "missing module",
				module:    "ctrsploit_e2e_missing_module_7f6b3c2a",
				satisfied: false,
			},
		},
	}
	testcases, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			satisfied, err := New(testcase.module).Check()
			require.NoError(t, err)
			require.Equal(t, testcase.satisfied, satisfied)
		})
	}
}
