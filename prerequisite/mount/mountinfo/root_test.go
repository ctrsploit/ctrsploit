package mountinfo

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_HostsRootContainsDocker(t *testing.T) {
	tests := map[string]struct {
		satisfied bool
	}{
		"docker-v28.3.2": {
			satisfied: true,
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		satisfied, err := HostsRootContainsDocker.Check()
		assert.NoError(t, err)
		assert.Equal(t, test.satisfied, satisfied)
	})
}

func TestE2E_HostnameRootContainsContainerd(t *testing.T) {
	tests := map[string]struct {
		satisfied bool
	}{
		"kubernetes-v1.33.1-calico": {
			satisfied: true,
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		satisfied, err := HostnameRootContainsContainerd.Check()
		assert.NoError(t, err)
		assert.Equal(t, test.satisfied, satisfied)
	})
}
