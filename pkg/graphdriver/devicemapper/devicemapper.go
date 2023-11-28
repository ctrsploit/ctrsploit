package devicemapper

// TODO: move to dir pkg/graphdriver/

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/module"
	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
	"regexp"
	"strings"
)

type DeviceMapper struct {
	AlreadyInit                      bool
	Loaded                           bool
	Used                             bool
	NumberOfDmUsedInRunningContainer int
	HostPath                         string
}

// https://tldp.org/HOWTO/LVM-HOWTO/builddmmod.html

const DirSysDeviceMapper = "/sys/class/misc/device-mapper"
const ProcDevices = "/proc/devices"
const BlockNameDeviceMapper = "device-mapper"
const ModuleNameThinPool = "dm_thin_pool"
const ModuleNameDm = "dm_mod"

func (d *DeviceMapper) Init() (err error) {
	// is enabled
	d.Loaded, err = d.IsEnabled()
	if err != nil {
		return
	}

	if d.Loaded {
		d.Used, err = d.IsUsed()
		if err != nil {
			return
		}
	}

	if d.Used {
		d.NumberOfDmUsedInRunningContainer, err = d.Number()
		if err != nil {
			return
		}
		d.HostPath, err = d.HostPathOfCtrRootfs()
		if err != nil {
			return
		}
	}

	d.AlreadyInit = true
	return
}

// IsEnabled
// https://tldp.org/HOWTO/LVM-HOWTO/builddmmod.html
func (d *DeviceMapper) IsEnabled() (enabled bool, err error) {
	if d.AlreadyInit {
		enabled = d.Loaded
		return
	}
	if _, err := os.Lstat(DirSysDeviceMapper); !os.IsNotExist(err) {
		enabled = true
		return enabled, err
	}
	content, err := os.ReadFile(ProcDevices)
	if err != nil {
		awesome_error.CheckErr(err)
		return enabled, err
	}
	if strings.Contains(string(content), BlockNameDeviceMapper) {
		enabled = true
	}
	return
}

// Number
// `cat /sys/module/dm_thin_pool/refcnt`==`docker ps |wc -l`
// https://docs.docker.com/storage/storagedriver/device-mapper-driver/
func (d *DeviceMapper) Number() (number int, err error) {
	if d.AlreadyInit {
		number = d.NumberOfDmUsedInRunningContainer
		return
	}
	if enabled, err := d.IsEnabled(); err != nil {
		return number, err
	} else {
		if enabled {
			number, err = module.RefCount(ModuleNameThinPool)
			if os.IsNotExist(err) {
				return number, nil
			}
			if err != nil {
				return number, err
			}
		}
	}
	return
}

func (d *DeviceMapper) IsUsed() (used bool, err error) {
	number, err := d.Number()
	if err != nil {
		return
	}
	if number > 0 {
		used = true
		return
	}
	if d.Loaded {
		number, err := module.RefCount(ModuleNameDm)
		if os.IsNotExist(err) {
			return used, nil
		}
		if err != nil {
			return used, err
		}
		used = number > 0
	}
	return
}

func (d *DeviceMapper) HostPathOfCtrRootfs() (host string, err error) {
	if d.AlreadyInit {
		host = d.HostPath
		return
	}
	used, err := d.IsUsed()
	if err != nil {
		return
	}
	if used {
		mount, err := mountinfo.RootMount()
		if err != nil {
			return host, err
		}
		// assert type
		if !mountinfo.IsDeviceMapper(mount) {
			awesome_error.CheckWarning(fmt.Errorf("not a device-mapper: %+v", mount))
			return host, err
		}
		pattern := regexp.MustCompile("-\\d+:\\d+-[0-9a-f]+-([0-9a-f]+)")
		matches := pattern.FindStringSubmatch(mount.Source)
		if len(matches) != 2 {
			err = fmt.Errorf("not a device-mapper: %+v", mount)
			awesome_error.CheckFatal(err)
		}
		dm := matches[1]
		host = fmt.Sprintf("/var/lib/docker/devicemapper/mnt/%s%s", dm, mount.Root)
	}
	return
}
