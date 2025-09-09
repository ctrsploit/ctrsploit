package selinux

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_Selinux(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	allTestcases := map[string][]struct {
		name            string
		kernelSupported bool
		enabled         bool
		privileged      bool
	}{
		"docker-v28.0.4_centos-stream9": {
			{
				name:            "privileged selinux",
				kernelSupported: true,
				enabled:         true,
				privileged:      true,
			},
		},
	}
	testcases, ok := allTestcases[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			kernelSupported, err := KernelSupported()
			assert.NoError(t, err)
			assert.Equal(t, testcase.kernelSupported, kernelSupported)
			assert.Equal(t, testcase.enabled, IsEnabled())
			assert.Equal(t, testcase.privileged, IsSelinuxPrivileged())
		})
	}
}
