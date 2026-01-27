package service_account

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_ServiceAccountMounted(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	allTestcases := map[string]struct {
		Satisfied bool
	}{
		"mounted": {
			Satisfied: true,
		},
		"not_mounted": {
			Satisfied: false,
		},
	}
	testcase, ok := allTestcases[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(testEnv, func(t *testing.T) {
		satisfied, err := Mounted.Check()
		assert.NoError(t, err)
		assert.Equal(t, testcase.Satisfied, satisfied)
	})
}
