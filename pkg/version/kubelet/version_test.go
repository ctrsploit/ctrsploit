package kubelet

import (
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"
)

func TestE2E_Versions(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	tests := map[string]struct {
		expected []*semver.Version
	}{
		"kubernetes-v1.34.0-host": {
			expected: []*semver.Version{semver.MustParse("v1.34.0")},
		},
	}
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(testEnv, func(t *testing.T) {
		versions, err := Versions()
		require.NoError(t, err)
		require.Equal(t, test.expected, versions)
	})
}

func TestE2E_VersionsByK8sApi(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	tests := map[string]struct {
		expected []*semver.Version
	}{
		"kubernetes-v1.34.0-host": {
			expected: []*semver.Version{semver.MustParse("v1.34.0")},
		},
	}
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(testEnv, func(t *testing.T) {
		versions, err := VersionsByK8sApi()
		require.NoError(t, err)
		require.Equal(t, test.expected, versions)
	})
}

func TestE2E_VersionByCli(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	tests := map[string]struct {
		expected *semver.Version
	}{
		"kubernetes-v1.34.0-host": {
			expected: semver.MustParse("v1.34.0"),
		},
	}
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(testEnv, func(t *testing.T) {
		version, err := VersionByCli()
		require.NoError(t, err)
		require.Equal(t, test.expected, version)
	})
}
