package block

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/google/cadvisor/fs"
	"github.com/google/cadvisor/machine"
	"github.com/google/cadvisor/utils/sysfs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetBlockDeviceInfo(t *testing.T) {
	info, err := GetBlockDeviceInfo()
	assert.NoError(t, err)
	spew.Dump(info)
}

func TestInfo(t *testing.T) {
	fsinfo, err := fs.NewFsInfo(fs.Context{})
	assert.NoError(t, err)
	info, err := machine.Info(sysfs.NewRealSysFs(), fsinfo, true)
	assert.NoError(t, err)
	spew.Dump(info)
}
