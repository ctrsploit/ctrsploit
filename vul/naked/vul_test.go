package naked

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_check(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	allTestcases := map[string][]struct {
		name  string
		naked bool
	}{
		"docker-v28.3.2_naked": {
			{
				name:  "naked",
				naked: true,
			},
		},
		"docker-v28.3.2_default": {
			{
				name:  "ubuntu default",
				naked: false,
			},
		},
	}
	testcases, ok := allTestcases[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			satisfied, err := Vul.CheckSecPrerequisites.Check()
			assert.NoError(t, err)
			assert.Equal(t, testcase.naked, satisfied)
		})
	}
}
