package storagedriver

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/hostpath/rootfs"
	"github.com/ctrsploit/ctrsploit/pkg/runtime"
	"github.com/ctrsploit/ctrsploit/pkg/storagedriver"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	spec "github.com/ctrsploit/sploit-spec/pkg/env/container/storagedriver"
)

const CommandName = "storage-driver"

func StorageDriver() (spec.StorageDriver, error) {
	driver, err := storagedriver.GetStorageDriver()
	if err != nil {
		return spec.StorageDriver{}, fmt.Errorf("could not get storage driver: %w", err)
	}
	enabled, err := driver.Enabled()
	if err != nil {
		return spec.StorageDriver{}, fmt.Errorf("failed to call driver.Enabled(): %w", err)
	}
	used, err := driver.Used()
	if err != nil {
		return spec.StorageDriver{}, fmt.Errorf("failed to call driver.Used(): %w", err)
	}
	number, err := driver.Number()
	if err != nil {
		return spec.StorageDriver{}, fmt.Errorf("failed to call driver.Number(): %w", err)
	}
	path, err := rootfs.HostPath(runtime.GetType(), driver.Type())
	if err != nil {
		return spec.StorageDriver{}, fmt.Errorf("failed to call rootfs.HostPath(): %w", err)
	}
	return spec.StorageDriver{
		Type:    driver.Type(),
		Enabled: enabled,
		Used:    used,
		Number:  number,
		Rootfs:  path,
	}, nil
}

func Filesystem() (filesystem container.Filesystem, err error) {
	driver, err := StorageDriver()
	return container.Filesystem{
		StorageDriver: driver,
	}, err
}
