package graphdriver

import (
	"ctrsploit/log"
	"ctrsploit/pkg/graphdriver/devicemapper"
	"ctrsploit/pkg/graphdriver/overlay"
	"ctrsploit/util"
	"fmt"
	"github.com/fatih/color"
)

const CommandName = "graphdriver"

func Overlay() (err error) {
	o := &overlay.Overlay{}
	err = o.Init()
	if err != nil {
		return
	}
	info := fmt.Sprintf("===========Overlay=========\nOverlay enabled: %v\n", util.ColorfulTickOrBallot(o.Loaded))
	if o.Loaded {
		info += fmt.Sprintf("Overlay used: %v\n", util.ColorfulTickOrBallot(o.Used))
		if o.Used {
			info += fmt.Sprintf("The number of overlayfs mounted: %v (equal to the number of containers)", color.HiGreenString(fmt.Sprintf("%d", o.Number)))
			if len(o.HostPath) > 0 {
				info += fmt.Sprintf("\nThe host path of container's rootfs: %s", color.HiGreenString(o.HostPath))
			}
			info += "\n"
		}
	}
	log.Logger.Info(info)
	return
}

func DeviceMapper() (err error) {
	d := &devicemapper.DeviceMapper{}
	err = d.Init()
	if err != nil {
		return
	}
	info := fmt.Sprintf("===========DeviceMapper=========\nDeviceMapper enabled: %v\n", util.ColorfulTickOrBallot(d.Loaded))
	if d.Loaded {
		info += fmt.Sprintf("DeviceMapper used: %v\n", util.ColorfulTickOrBallot(d.Used))
		if d.Used {
			info += fmt.Sprintf("The number of devicemapper used in running container: %v", color.HiGreenString(fmt.Sprintf("%d", d.NumberOfDmUsedInRunningContainer)))
			if d.NumberOfDmUsedInRunningContainer > 0 {
				info += " ( =(count(running containers)+1) )"
				if len(d.HostPath) > 0 {
					info += fmt.Sprintf("\nThe host path of container's rootfs: %s", color.HiGreenString(d.HostPath))
				}
			}
			info += "\n"
		}
	}
	log.Logger.Info(info)
	return
}
