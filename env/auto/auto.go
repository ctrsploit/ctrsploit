package auto

import (
	"github.com/ctrsploit/ctrsploit/env/apparmor"
	"github.com/ctrsploit/ctrsploit/env/capability"
	"github.com/ctrsploit/ctrsploit/env/cgroups"
	"github.com/ctrsploit/ctrsploit/env/mountinfo"
	"github.com/ctrsploit/ctrsploit/env/namespace"
	"github.com/ctrsploit/ctrsploit/env/seccomp"
	"github.com/ctrsploit/ctrsploit/env/selinux"
	"github.com/ctrsploit/ctrsploit/env/storagedriver"
	"github.com/ctrsploit/ctrsploit/env/where"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
)

func Basic() (basic container.Basic, err error) {
	w, err := where.Where()
	if err != nil {
		return
	}
	m, err := mountinfo.MountInfo()
	if err != nil {
		return
	}
	basic = container.Basic{
		Where:         w,
		KernelVersion: "", //TODO
		MountInfo:     m,
	}
	return
}

func LinuxSecurityFeature() (lsf container.LinuxSecurityFeature, err error) {
	cap, err := capability.Capability()
	if err != nil {
		return
	}
	aa, err := apparmor.Apparmor()
	if err != nil {
		return
	}
	se, err := selinux.Selinux()
	if err != nil {
		return
	}
	sc, err := seccomp.Seccomp()
	if err != nil {
		return
	}
	ns, err := namespace.Namespace()
	if err != nil {
		return
	}
	cg, err := cgroups.Cgroups()
	if err != nil {
		return
	}
	fs, err := storagedriver.Filesystem()
	if err != nil {
		return
	}
	lsf = container.LinuxSecurityFeature{
		Credential:   container.Credential{}, //TODO
		Capabilities: cap,
		LSM: container.LSM{
			Apparmor: aa,
			SELinux:  se,
		},
		Seccomp:    sc,
		Namespace:  ns,
		CGroups:    cg,
		Filesystem: fs,
	}
	return
}

func Auto() (env container.Env, err error) {
	basic, err := Basic()
	if err != nil {
		return
	}
	lsf, err := LinuxSecurityFeature()
	if err != nil {
		return
	}
	env = container.Env{
		Basic:                basic,
		LinuxSecurityFeature: lsf,
		Cluster:              container.Cluster{}, //TODO
		Advance:              container.Advance{}, //TODO
	}
	return
}
