package namespace

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListNamespaceDir(t *testing.T) {
	namespaces, err := ListNamespaceDir("/proc/self/ns")
	assert.NoError(t, err)
	assert.True(t, len(namespaces) > 0)
	for _, ns := range []string{
		"ipc", "mnt", "net", "pid", "user", "uts",
	} {
		assert.True(t, namespaces[ns] > 0)
	}
}

func TestReadInodeNumberUnderProc(t *testing.T) {
	inoMap, err := ReadInodeNumberMapUnderProc("/proc")
	assert.NoError(t, err)
	for _, ino := range inoMap {
		assert.True(t, ino > 0xF0000000)
		assert.True(t, ino < 0xFFFFFFF0)
	}
}
