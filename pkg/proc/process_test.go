package proc

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsSheBang(t *testing.T) {
	shebang := filepath.Join(t.TempDir(), "ctrsploit_shebang")
	assert.NoError(t, os.WriteFile(shebang, []byte("#!/bin/bash\nsleep 10"), 0o755))
	cmd := exec.Command("/bin/bash", "-c", shebang)
	assert.NoError(t, cmd.Start())
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	})
	time.Sleep(time.Second)
	isSheBang, err := IsSheBang(cmd.Process.Pid)
	assert.NoError(t, err)
	assert.True(t, isSheBang)
}
