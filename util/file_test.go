package util

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestCheckFileExists(t *testing.T) {
	assert.True(t, CheckPathExists("/etc/passwd"))
	assert.False(t, CheckPathExists("/not_exists"))
}

func TestReadIntFromFile(t *testing.T) {
	result, err := ReadIntFromFile("/proc/sys/kernel/pid_max")
	assert.NoError(t, err)
	assert.True(t, result > 0)
}

func TestReplaceContent(t *testing.T) {
	assert.NoError(t, ioutil.WriteFile("/tmp/replace_test", []byte("source"), 0755))
	assert.NoError(t, ReplaceContent("/tmp/replace_test", []byte("source"), []byte("dest")))
	content, err := ioutil.ReadFile("/tmp/replace_test")
	assert.NoError(t, err)
	assert.Equal(t, []byte("dest"), content)
}
