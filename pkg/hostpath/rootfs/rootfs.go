package rootfs

import (
	"fmt"

	"github.com/ctrsploit/ctrsploit/pkg/runtime"
	"github.com/ctrsploit/sploit-spec/pkg/env/container/storagedriver"
)

type Driver interface {
	Is() (bool, error)
	RootPath() (string, error)
}

const (
	ErrorUnsupportedRuntime = "unsupported runtime type %s to get root path"
	ErrorUnsupportedDriver  = "unsupported storage driver %s for runtime %s to get root path"
)

func HostPath(runtimeType runtime.Type, storageDriverType storagedriver.Type) (string, error) {
	var driver Driver
	if runtimeType&runtime.TypeDocker != 0 {
		switch storageDriverType {
		case storagedriver.TypeOverlay:
			driver = DockerOverlay{}
		case storagedriver.TypeDeviceMapper:
			driver = DockerDeviceMapper{}
		default:
			return "", fmt.Errorf(ErrorUnsupportedDriver, storageDriverType, runtimeType)
		}
	} else if runtimeType&runtime.TypeContainerd != 0 {
		switch storageDriverType {
		case storagedriver.TypeOverlay:
			driver = ContainerdOverlay{}
		default:
			return "", fmt.Errorf(ErrorUnsupportedDriver, storageDriverType, runtimeType)
		}
	} else {
		return "", fmt.Errorf(ErrorUnsupportedRuntime, runtimeType)
	}
	return driver.RootPath()
}
