package block

import (
	info "github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/utils/sysfs"
	"github.com/google/cadvisor/utils/sysinfo"
)

func GetBlockDeviceInfo() (map[string]info.DiskInfo, error) {
	return sysinfo.GetBlockDeviceInfo(sysfs.NewRealSysFs())
}
