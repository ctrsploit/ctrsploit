package prerequisite

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/pkg/kernel/uname"
)

type KernelVersion struct {
	ExpectedMinVersion string
	ExpectedMaxVersion string
	BasePrerequisite
}

var (
	KernelSupportsCgroupNamespace = KernelVersion{
		ExpectedMinVersion: "4.6",
		ExpectedMaxVersion: "",
		BasePrerequisite: BasePrerequisite{
			Name: "Kernel Supports Cgroup namespace",
			Info: "kernel version >= v4.6",
		},
	}
	KernelSupportsTimeNamespace = KernelVersion{
		ExpectedMinVersion: "5.6",
		ExpectedMaxVersion: "",
		BasePrerequisite: BasePrerequisite{
			Name: "Kernel Supports Time namespace",
			Info: "kernel version >= v5.6",
		},
	}
)

func (p *KernelVersion) check(version string) (satisfied bool) {
	satisfied = true
	if p.ExpectedMinVersion != "" {
		satisfied = satisfied && (p.ExpectedMinVersion < version || uname.VersionEqual(p.ExpectedMinVersion, version))
	}
	if p.ExpectedMaxVersion != "" {
		satisfied = satisfied && (p.ExpectedMaxVersion > version || uname.VersionEqual(p.ExpectedMaxVersion, version))
	}
	log.Logger.Debugf("%s <= %s <= %s: %t\n", p.ExpectedMinVersion, version, p.ExpectedMaxVersion, satisfied)
	return
}

func (p *KernelVersion) Check() (err error) {
	version, err := uname.Release()
	if err != nil {
		return
	}
	p.Satisfied = p.check(version)
	return
}
