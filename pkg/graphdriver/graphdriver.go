package graphdriver

import (
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver/devicemapper"
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver/overlay"
	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
)

type Interface interface {
	Init() (err error)
	IsEnabled() (enabled bool, err error)
	IsUsed() (used bool, err error)
	HostPathOfCtrRootfs() (host string, err error)
	Number() (number int, err error)
}

type Type int

const (
	Unkown Type = iota
	Overlay
	DeviceMapper
)

type GraphDriver struct {
	Type         Type
	Rootfs       string
	overlay      overlay.Overlay
	deviceMapper devicemapper.DeviceMapper
}

func (g *GraphDriver) DetectType() (err error) {
	mount, err := mountinfo.RootMount()
	if err != nil {
		return
	}
	if mountinfo.IsOverlay(mount) {
		g.Type = Overlay
	} else if mountinfo.IsDeviceMapper(mount) {
		g.Type = DeviceMapper
	}
	return
}

func (g *GraphDriver) Init() (err error) {
	err = g.DetectType()
	if err != nil {
		return
	}
	switch g.Type {
	case Overlay:
		err = g.overlay.Init()
		if err != nil {
			return
		}
		g.Rootfs = g.overlay.HostPath
	case DeviceMapper:
		err = g.deviceMapper.Init()
		if err != nil {
			return
		}
		g.Rootfs = g.deviceMapper.HostPath
	}
	return
}
