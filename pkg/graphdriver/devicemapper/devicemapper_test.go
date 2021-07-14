package devicemapper

import (
	"github.com/ssst0n3/awesome_libs/log"
	"testing"
)

func TestDeviceMapper_Number(t *testing.T) {
	d := &DeviceMapper{}
	log.Logger.Info(d.Number())
}
