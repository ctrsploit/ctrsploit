package syscall

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_Syscall_Check(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	allTestcases := map[string][]struct {
		available bool
	}{
		"seccomp-enabled": {
			{
				available: false,
			},
		},
		"seccomp-unconfined": {
			{
				available: true,
			},
		},
	}
	testcases, ok := allTestcases[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	for _, testcase := range testcases {
		t.Run(testEnv, func(t *testing.T) {
			satisfied, err := Unshare.Check()
			assert.NoError(t, err)
			assert.Equal(t, testcase.available, satisfied)
		})
	}
}
