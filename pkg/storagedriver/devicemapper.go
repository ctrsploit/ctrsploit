package storagedriver

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ctrsploit/ctrsploit/pkg/module"
	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ctrsploit/sploit-spec/pkg/env/container/storagedriver"
)

// https://tldp.org/HOWTO/LVM-HOWTO/builddmmod.html

const DirSysDeviceMapper = "/sys/class/misc/device-mapper"
const ProcDevices = "/proc/devices"
const BlockNameDeviceMapper = "device-mapper"
const ModuleNameThinPool = "dm_thin_pool"

type DeviceMapper struct {
}

func NewDeviceMapper() *DeviceMapper {
	return &DeviceMapper{}
}

func (d *DeviceMapper) Type() storagedriver.Type {
	return storagedriver.TypeDeviceMapper
}

// Enabled
// https://tldp.org/HOWTO/LVM-HOWTO/builddmmod.html
func (d *DeviceMapper) Enabled() (bool, error) {
	exists, _ := internal.CheckPathExists(DirSysDeviceMapper)
	if exists {
		return true, nil
	}
	content, err := os.ReadFile(ProcDevices)
	if err != nil {
		return false, fmt.Errorf("error reading %s: %w", ProcDevices, err)
	}
	return strings.Contains(string(content), BlockNameDeviceMapper), nil
}

// Number
// `cat /sys/module/dm_thin_pool/refcnt`==`docker ps |wc -l`
// https://docs.docker.com/storage/storagedriver/device-mapper-driver/
func (d *DeviceMapper) Number() (int, error) {
	number, err := module.RefCount(ModuleNameThinPool)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		} else {
			return 0, fmt.Errorf("error reading %s: %w", ModuleNameThinPool, err)
		}
	}
	return number, nil
}

func (o *DeviceMapper) Used() (bool, error) {
	var errs []error
	info, err := mountinfo.RootMount()
	if err != nil {
		errs = append(errs, fmt.Errorf("error getting root mount info: %w", err))
	} else if mountinfo.IsDeviceMapper(info) {
		return true, nil
	}
	number, err := o.Number()
	if err != nil {
		errs = append(errs, fmt.Errorf("error getting overlay number: %w", err))
	} else if number > 0 {
		return true, nil
	}
	return false, errors.Join(errs...)
}
