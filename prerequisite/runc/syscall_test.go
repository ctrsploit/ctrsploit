package runc

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_EnsureCloned(t *testing.T) {
	tests := map[string]struct {
		satisfied bool
	}{
		"runc-v0.0.1": {
			satisfied: false,
		},
		"runc-v0.1.0": {
			satisfied: false,
		},
		"runc-v1.0.0-rc6": {
			satisfied: false,
		},
		"runc-v1.0.0-rc7": {
			satisfied: true,
		},
		"runc-v1.0.0": {
			satisfied: true,
		},
		"runc-v1.1.0": {
			satisfied: true,
		},
		"runc-v1.1.15": {
			satisfied: true,
		},
		"runc-v1.2.0-rc.1": {
			satisfied: false,
		},
		"runc-v1.3.0": {
			satisfied: false,
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		satisfied, err := EnsureCloned.Check()
		assert.NoError(t, err)
		assert.Equal(t, test.satisfied, satisfied)
	})
}
