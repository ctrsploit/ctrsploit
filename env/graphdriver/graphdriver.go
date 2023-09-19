package graphdriver

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver/devicemapper"
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver/overlay"
	"github.com/ctrsploit/ctrsploit/util"
	"github.com/fatih/color"
)

const CommandGraphdriverName = "graphdriver"

func Overlay() (err error) {
	o := &overlay.Overlay{}
	err = o.Init()
	if err != nil {
		return
	}
	info := fmt.Sprintf("===========Overlay=========\nOverlay enabled: %v", util.ColorfulTickOrBallot(o.Loaded))
	if o.Loaded {
		info += fmt.Sprintf("\nOverlay used: %v", util.ColorfulTickOrBallot(o.Used))
		if o.Used {
			info += fmt.Sprintf("\nThe number of overlayfs mounted: %v (equal to the number of containers)", color.HiGreenString(fmt.Sprintf("%d", o.Number)))
			if len(o.HostPath) > 0 {
				info += fmt.Sprintf("\nThe host path of container's rootfs: %s", color.HiGreenString(o.HostPath))
			}
		}
	}
	fmt.Printf("%s\n\n", info)
	return
}

func DeviceMapper() (err error) {
	d := &devicemapper.DeviceMapper{}
	err = d.Init()
	if err != nil {
		return
	}
	info := fmt.Sprintf("===========DeviceMapper=========\nDeviceMapper enabled: %v", util.ColorfulTickOrBallot(d.Loaded))
	if d.Loaded {
		info += fmt.Sprintf("\nDeviceMapper used: %v", util.ColorfulTickOrBallot(d.Used))
		if d.Used {
			info += fmt.Sprintf("\nThe number of devicemapper used in running container: %v", color.HiGreenString(fmt.Sprintf("%d", d.NumberOfDmUsedInRunningContainer)))
			if d.NumberOfDmUsedInRunningContainer > 0 {
				info += " ( =(count(running containers)+1) )"
				if len(d.HostPath) > 0 {
					info += fmt.Sprintf("\nThe host path of container's rootfs: %s", color.HiGreenString(d.HostPath))
				}
			}
		}
	}
	fmt.Printf("%s\n\n", info)
	return
}
