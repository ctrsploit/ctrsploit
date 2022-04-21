package util

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestGetProcessNameByPid(t *testing.T) {
	expectedProcess := "cat"
	cmd := exec.Command(expectedProcess)
	assert.NoError(t, cmd.Start())
	processName, err := GetProcessNameByPid(cmd.Process.Pid)
	assert.NoError(t, err)
	assert.Equal(t, expectedProcess, processName)
}

func TestIsSheBang(t *testing.T) {
	shebang := "/tmp/ctrsploit_shebang"
	assert.NoError(t, os.Remove(shebang))
	assert.NoError(t, ioutil.WriteFile(shebang, []byte("#!/bin/bash\nsleep 10"), 0755))
	cmd := exec.Command("/bin/bash", "-c", shebang)
	assert.NoError(t, cmd.Start())
	time.Sleep(time.Second)
	isSheBang, err := IsSheBang(cmd.Process.Pid)
	assert.NoError(t, err)
	assert.True(t, isSheBang)
}
