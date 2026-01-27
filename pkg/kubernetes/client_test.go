package kubernetes

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_GetKubernetesConfig(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	tests := map[string]struct {
		canGetKubeconfig bool
	}{
		"kubernetes-v1.34.0-host": {
			canGetKubeconfig: true,
		},
	}
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(testEnv, func(t *testing.T) {
		config, err := GetKubernetesConfig()
		assert.Equal(t, test.canGetKubeconfig, err == nil, fmt.Sprintf("canGetKubeconfig: %v, got error: %v", test.canGetKubeconfig, err))
		assert.Equal(t, test.canGetKubeconfig, config != nil, fmt.Sprintf("canGetKubeconfig: %v, got config: %v", test.canGetKubeconfig, config))
	})
}

func TestE2E_GetKubernetesClient(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	tests := map[string]struct {
		canGetClient bool
	}{
		"kubernetes-v1.34.0-host": {
			canGetClient: true,
		},
	}
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(testEnv, func(t *testing.T) {
		cli, err := GetKubernetesClient()
		assert.Equal(t, test.canGetClient, err == nil, fmt.Sprintf("canGetKubeconfig: %v, got error: %v", test.canGetClient, err))
		assert.Equal(t, test.canGetClient, cli != nil, fmt.Sprintf("canGetKubeconfig: %v, got config: %v", test.canGetClient, cli))
	})
}
