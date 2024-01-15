package pkg

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/version/libseccomp"
	"github.com/ctrsploit/ctrsploit/pkg/version/version"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestElfStripped(t *testing.T) {
	resp, err := http.DefaultClient.Get("https://github.com/opencontainers/runc/releases/download/v1.1.11/runc.amd64")
	assert.NoError(t, err)
	content, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	path := fmt.Sprintf("/tmp/runc-%s", uuid.New())
	err = os.WriteFile(path, content, 0644)
	assert.NoError(t, err)
	stripped, err := ElfStripped(path)
	assert.True(t, stripped)
	assert.NoError(t, err)
	assert.NoError(t, os.Remove(path))
}

func TestExecuteRuncVersion(t *testing.T) {
	resp, err := http.DefaultClient.Get("https://github.com/opencontainers/runc/releases/download/v1.1.11/runc.amd64")
	assert.NoError(t, err)
	content, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	path := fmt.Sprintf("/tmp/runc-%s", uuid.New())
	err = os.WriteFile(path, content, 0644)
	assert.NoError(t, err)
	ver, err := ExecuteRuncVersion(path)
	assert.NoError(t, err)
	assert.Equal(t, libseccomp.Version{
		Number: version.Number{
			Major: 2,
			Minor: 5,
			Patch: 4,
			Rc:    -1,
			Beta:  -1,
			Init:  true,
		},
	}, ver)
	assert.NoError(t, os.Remove(path))
}

func TestReadSymbol(t *testing.T) {
	resp, err := http.DefaultClient.Get("https://github.com/opencontainers/runc/releases/download/v1.0.0-rc92/runc.amd64")
	assert.NoError(t, err)
	content, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	path := fmt.Sprintf("/tmp/runc-%s", uuid.New())
	err = os.WriteFile(path, content, 0644)
	assert.NoError(t, err)
	ver, err := ReadSymbol(path)
	assert.NoError(t, err)
	assert.Equal(t, libseccomp.Version{
		Number: version.Number{
			Major: 2,
			Minor: 4,
			Patch: 1,
			Rc:    -1,
			Beta:  -1,
			Init:  true,
		},
	}, ver)
	assert.NoError(t, os.Remove(path))
}
