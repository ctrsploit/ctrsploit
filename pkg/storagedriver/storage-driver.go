package storagedriver

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ctrsploit/sploit-spec/pkg/env/container/storagedriver"
)

type Driver interface {
	Type() storagedriver.Type
	Enabled() (bool, error)
	Used() (bool, error)
	Number() (int, error)
}

func GetStorageDriver() (Driver, error) {
	t, err := GetType()
	if err != nil {
		return nil, err
	}
	switch t {
	case storagedriver.TypeOverlay:
		return NewOverlay(), nil
	case storagedriver.TypeDeviceMapper:
		return NewDeviceMapper(), nil
	default:
		return nil, fmt.Errorf("unknown storage driver: %s", t)
	}
}

func GetType() (storagedriver.Type, error) {
	info, err := mountinfo.RootMount()
	if err != nil {
		return storagedriver.TypeUnknown, fmt.Errorf("error getting root's mount info: %w", err)
	}
	if mountinfo.IsDeviceMapper(info) {
		return storagedriver.TypeDeviceMapper, nil
	}
	if mountinfo.IsOverlay(info) {
		return storagedriver.TypeOverlay, nil
	}
	return storagedriver.TypeUnknown, fmt.Errorf("unknown storage type")
}
