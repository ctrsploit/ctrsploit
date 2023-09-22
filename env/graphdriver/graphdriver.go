package graphdriver

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/internal/colorful"
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver/devicemapper"
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver/overlay"
	"github.com/ssst0n3/awesome_libs"
)

const CommandGraphdriverName = "graphdriver"

var (
	tplOverlay = `
===========Overlay===========
{.enabled}  Enabled
{.used}`
	tplOverlayUsed = `{.used}  Used
{.details}`
	tplOverlayDetails = `The number of overlayfs mounted(equal to the number of containers)
{.number}
The host path of container's rootfs:
{.host}
`
)

func Overlay() (err error) {
	o := &overlay.Overlay{}
	err = o.Init()
	if err != nil {
		return
	}
	used := ""
	if o.Loaded {
		details := ""
		if o.Used {
			details = awesome_libs.Format(tplOverlayDetails, awesome_libs.Dict{
				"number": colorful.Safe(fmt.Sprintf("%d", o.Number)),
				"host":   colorful.Safe(o.HostPath),
			})
		}
		used = awesome_libs.Format(tplOverlayUsed, awesome_libs.Dict{
			"used":    colorful.TickOrBallot(o.Used),
			"details": details,
		})
	}
	info := awesome_libs.Format(tplOverlay, awesome_libs.Dict{
		"enabled": colorful.TickOrBallot(o.Loaded),
		"used":    used,
	})
	print(info)
	return
}

var (
	tplDeviceMapper = `
===========DeviceMapper===========
{.enabled}  Enabled
{.used}`
	tplDeviceMapperUsed = `{.used}  Used
{.details}`
	tplDeviceMapperDetails = `The number of devicemapper used in running container:
{.number}
The host path of container's rootfs:
{.host}
`
)

func DeviceMapper() (err error) {
	d := &devicemapper.DeviceMapper{}
	err = d.Init()
	if err != nil {
		return
	}
	used := ""
	if d.Loaded {
		details := ""
		if d.Used {
			details = awesome_libs.Format(tplOverlayDetails, awesome_libs.Dict{
				"number": colorful.Safe(fmt.Sprintf("%d", d.NumberOfDmUsedInRunningContainer)),
				"host":   colorful.Safe(d.HostPath),
			})
		}
		used = awesome_libs.Format(tplDeviceMapperUsed, awesome_libs.Dict{
			"used":    colorful.TickOrBallot(d.Used),
			"details": details,
		})
	}
	info := awesome_libs.Format(tplDeviceMapper, awesome_libs.Dict{
		"enabled": colorful.TickOrBallot(d.Loaded),
		"used":    used,
	})
	print(info)
	return
}
