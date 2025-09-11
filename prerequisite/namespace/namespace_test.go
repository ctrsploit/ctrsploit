package namespace

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_Namespace_Check(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	allTestcases := map[string][]struct {
		host bool
	}{
		"pid-host": {
			{
				host: true,
			},
		},
		"pid-default": {
			{
				host: false,
			},
		},
	}
	testcases, ok := allTestcases[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	for _, testcase := range testcases {
		t.Run(testEnv, func(t *testing.T) {
			satisfied, err := PidNamespaceLevelHost.Check()
			assert.NoError(t, err)
			assert.Equal(t, testcase.host, satisfied)
		})
	}
}
