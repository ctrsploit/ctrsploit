package hostname

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_K8sHostnameMatch(t *testing.T) {
	tests := map[string]struct {
		satisfied bool
	}{
		"docker-v28.3.2": {
			satisfied: false,
		},
		"nerdctl-v2.1.2": {
			satisfied: false,
		},
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
		satisfied, err := K8sDeploymentHostname.Check()
		assert.NoError(t, err)
		assert.Equal(t, test.satisfied, satisfied)
	})
}
