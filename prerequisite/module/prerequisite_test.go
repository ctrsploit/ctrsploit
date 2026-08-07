package module

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestE2E_ModulePrerequisites(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	tests := map[string][]struct {
		name               string
		module             string
		loadedSatisfied    bool
		availableSatisfied bool
	}{
		"docker-v28.3.2": {
			{
				name:               "loaded overlay module",
				module:             "overlay",
				loadedSatisfied:    true,
				availableSatisfied: true,
			},
			{
				name:               "missing module",
				module:             "ctrsploit_e2e_missing_module_7f6b3c2a",
				loadedSatisfied:    false,
				availableSatisfied: false,
			},
		},
	}
	testcases, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}

	for _, testcase := range testcases {
		t.Run(testcase.name+"/loaded", func(t *testing.T) {
			satisfied, err := New(testcase.module).Check()
			require.NoError(t, err)
			require.Equal(t, testcase.loadedSatisfied, satisfied)
		})
		t.Run(testcase.name+"/available", func(t *testing.T) {
			satisfied, err := NewAvailable(testcase.module).Check()
			require.NoError(t, err)
			require.Equal(t, testcase.availableSatisfied, satisfied)
		})
	}
}
