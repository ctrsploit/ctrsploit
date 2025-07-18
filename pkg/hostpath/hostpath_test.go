package hostpath

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"regexp"
	"testing"
)

func TestWritableAccessible(t *testing.T) {
	paths, err := WritableAccessible()
	assert.NoError(t, err)
	spew.Dump(paths)
}

func TestE2E_WritableAccessible(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	type expected struct {
		path        Path
		hostPattern string
	}

	paths, err := WritableAccessible()
	require.NoError(t, err)

	envExpectations := map[string][]expected{
		"docker-v19.03.13": {
			{
				path: Path{
					ContainerPath: "/",
					Type:          TypeRootfs,
				},
				hostPattern: `/var/lib/docker/overlay2/[0-9a-f]+/merged`,
			},
			{
				path: Path{
					ContainerPath: "/usr/bin/ctrsploit.test",
					Type:          TypeUserCustomBindMount,
					HostPath:      "/ctrsploit.test",
				},
			},
			{
				path: Path{
					ContainerPath: "/etc/hosts",
					Type:          TypeNetworkFiles,
				},
				hostPattern: `/var/lib/docker/containers/[0-9a-f]+/hosts`,
			},
			{
				path: Path{
					ContainerPath: "/etc/hostname",
					Type:          TypeNetworkFiles,
				},
				hostPattern: `/var/lib/docker/containers/[0-9a-f]+/hostname`,
			},
			{
				path: Path{
					ContainerPath: "/etc/resolv.conf",
					Type:          TypeNetworkFiles,
				},
				hostPattern: `/var/lib/docker/containers/[0-9a-f]+/resolv.conf`,
			},
		},
		"example": {
			{
				path: Path{
					ContainerPath: "/etc/resolv.conf",
					Type:          TypeNetworkFiles,
				},
				hostPattern: `/var/lib/docker/containers/[0-9a-f]+/resolv.conf`,
			},
		},
	}

	expectedTests, ok := envExpectations[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}

	assert.Equal(t, len(expectedTests), len(paths), "there're more than expected paths")

	actualPaths := make(map[string]Path)
	for _, p := range paths {
		actualPaths[p.ContainerPath] = p
	}

	for _, exp := range expectedTests {
		// Capture the current loop variable 'exp' to avoid issues with closure capturing.
		// In Go, the loop variable 'exp' is reused during each iteration.
		// Without this assignment, the closure would always reference the final value of 'exp'.
		exp := exp
		name := fmt.Sprintf("%s:%s", testEnv, exp.path.ContainerPath)
		t.Run(name, func(t *testing.T) {
			actual, exists := actualPaths[exp.path.ContainerPath]
			require.True(t, exists, "cannot find container path: %s", exp.path.ContainerPath)

			assert.Equal(t, exp.path.Type, actual.Type, "type not match: %s", exp.path.Type)
			if exp.hostPattern != "" {
				matched, err := regexp.MatchString(exp.hostPattern, actual.HostPath)
				require.NoError(t, err, "regex compiled failed: %s", exp.hostPattern)
				assert.True(t, matched, "host not match: %s", exp.path.ContainerPath)
			} else {
				assert.Equal(t, exp.path.HostPath, actual.HostPath, "host not match: %s", exp.path.HostPath)
			}
		})
	}
}
