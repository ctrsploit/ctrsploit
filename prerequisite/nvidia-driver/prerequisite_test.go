package nvidia_driver

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeviceNode(t *testing.T) {
	satisfied, err := newDeviceNode("/dev/null").Check()
	require.NoError(t, err)
	require.True(t, satisfied)

	satisfied, err = newDeviceNode(filepath.Join(t.TempDir(), "missing")).Check()
	require.NoError(t, err)
	require.False(t, satisfied)
}

func TestDeviceNodes(t *testing.T) {
	satisfied, err := newDeviceNodes("/dev/null", "/dev/null").Check()
	require.NoError(t, err)
	require.True(t, satisfied)

	satisfied, err = newDeviceNodes("/dev/null", filepath.Join(t.TempDir(), "missing")).Check()
	require.NoError(t, err)
	require.False(t, satisfied)
}

func TestVersionCheckPortable(t *testing.T) {
	// The host may be macOS without /proc or an NVIDIA GPU. The check must
	// still return one of the three verdicts without error.
	_, err := newVersion().Check()
	require.NoError(t, err)
}

func TestVerdictString(t *testing.T) {
	require.Equal(t, "vulnerable", VerdictVulnerable.String())
	require.Equal(t, "fixed", VerdictFixed.String())
	require.Equal(t, "unknown", VerdictUnknown.String())
}
