package source

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_RootMountInfoSourceContainsDocker(t *testing.T) {
	tests := map[string]struct {
		satisfied bool
	}{
		"docker-v23.0.0-devicemapper": {
			satisfied: true,
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		satisfied, err := RootMountInfoSourceContainsDocker.Check()
		assert.NoError(t, err)
		assert.Equal(t, test.satisfied, satisfied)
	})
}

func TestE2E_RootMountInfoVFSOptionsContainsDocker(t *testing.T) {
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
		satisfied, err := RootMountInfoVFSOptionsContainsDocker.Check()
		assert.NoError(t, err)
		assert.Equal(t, test.satisfied, satisfied)
	})
}

func TestE2E_RootMountInfoVFSOptionsContainsContainerd(t *testing.T) {
	tests := map[string]struct {
		satisfied bool
	}{
		"nerdctl-v2.1.2": {
			satisfied: true,
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		satisfied, err := RootMountInfoVFSOptionsContainsContainerd.Check()
		assert.NoError(t, err)
		assert.Equal(t, test.satisfied, satisfied)
	})
}
