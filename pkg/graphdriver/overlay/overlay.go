package overlay

import (
	"errors"
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/module"
	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
	"regexp"
)

// TODO: move to dir pkg/graphdriver/

type Overlay struct {
	AlreadyInit bool
	Loaded      bool
	number      int  // the number of overlays mounted
	Used        bool // the host path of container's rootfs
	HostPath    string
}

// Init Check Overlay is enabled and used
func (o *Overlay) Init() (err error) {
	o.Loaded, err = o.IsEnabled()
	if err != nil {
		return
	}
	if o.Loaded {
		o.Used, err = o.IsUsed()
		if err != nil {
			return err
		}
	}
	if o.Used {
		o.HostPath, err = o.HostPathOfCtrRootfs()
		if err != nil {
			return err
		}
	}
	o.AlreadyInit = true
	return
}

func (o *Overlay) IsEnabled() (enabled bool, err error) {
	if o.AlreadyInit {
		enabled = o.Loaded
		return
	}
	o.Number, err = module.RefCount("overlay")
	if err != nil {
		if os.IsNotExist(err) {
			enabled = false
			err = nil
		} else {
			awesome_error.CheckErr(err)
		}
		return
	}
	enabled = o.Number >= 0
	return
}

func (o *Overlay) IsUsed() (used bool, err error) {
	if o.AlreadyInit {
		used = o.Used
		return
	}
	loaded, err := o.IsEnabled()
	if err != nil {
		return
	}
	if loaded {
		o.number, err = module.RefCount("overlay")
		if err != nil {
			if os.IsNotExist(err) {
				mount, e := mountinfo.RootMount()
				err = e
				if err != nil {
					return
				}
				used = mountinfo.IsOverlay(mount)
			} else {
				awesome_error.CheckErr(err)
			}
			return
		}
		used = o.number >= 0
	}
	return
}

func (o *Overlay) HostPathOfCtrRootfs() (host string, err error) {
	if o.AlreadyInit {
		host = o.HostPath
		return
	}
	used, err := o.IsUsed()
	if err != nil {
		return
	}
	if used {
		mount, err := mountinfo.RootMount()
		if err != nil {
			return host, err
		}
		// assert type
		if !mountinfo.IsOverlay(mount) {
			err = errors.New(fmt.Sprintf("not a overlay, or you are not in the container: %+v", mount))
			awesome_error.CheckWarning(err)
			return host, err
		}
		pattern := regexp.MustCompile(",upperdir=(.*)/diff,")
		matches := pattern.FindStringSubmatch(mount.VFSOptions)
		if len(matches) != 2 {
			err = errors.New(fmt.Sprintf("Unkown VFSOptions: %+v, please add a issue to tell us", mount))
			awesome_error.CheckErr(err)
			return host, err
		}
		host = matches[1] + "/merged"
	}
	return
}

func (o *Overlay) Number() (number int, err error) {
	if o.AlreadyInit {
		number = o.number
		return
	}
	o.number, err = module.RefCount("overlay")
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		} else {
			awesome_error.CheckErr(err)
		}
		return
	}
	number = o.number
	return
}
