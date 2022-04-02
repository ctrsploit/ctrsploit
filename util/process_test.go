package util

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func TestGetProcessNameByPid(t *testing.T) {
	expectedProcess := "cat"
	cmd := exec.Command(expectedProcess)
	assert.NoError(t, cmd.Start())
	processName, err := GetProcessNameByPid(cmd.Process.Pid)
	assert.NoError(t, err)
	assert.Equal(t, expectedProcess, processName)
}
