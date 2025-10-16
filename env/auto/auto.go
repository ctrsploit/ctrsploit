package auto

import (
	"errors"

	"github.com/ctrsploit/ctrsploit/env/apparmor"
	"github.com/ctrsploit/ctrsploit/env/capability"
	"github.com/ctrsploit/ctrsploit/env/cgroups"
	"github.com/ctrsploit/ctrsploit/env/kernel"
	"github.com/ctrsploit/ctrsploit/env/mountinfo"
	"github.com/ctrsploit/ctrsploit/env/namespace"
	"github.com/ctrsploit/ctrsploit/env/seccomp"
	"github.com/ctrsploit/ctrsploit/env/selinux"
	"github.com/ctrsploit/ctrsploit/env/storagedriver"
	"github.com/ctrsploit/ctrsploit/env/where"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
)

func Basic() (container.Basic, error) {
	var errs []error
	w, err := where.Where()
	if err != nil {
		errs = append(errs, err)
	}
	m, err := mountinfo.MountInfo()
	if err != nil {
		errs = append(errs, err)
	}
	basic := container.Basic{
		Where:     w,
		MountInfo: m,
	}
	return basic, errors.Join(errs...)
}

func LinuxSecurityFeature() (container.LinuxSecurityFeature, error) {
	var errs []error
	caps, err := capability.Capability()
	if err != nil {
		errs = append(errs, err)
	}
	aa, err := apparmor.Apparmor()
	if err != nil {
		errs = append(errs, err)
	}
	se, err := selinux.Selinux()
	if err != nil {
		errs = append(errs, err)
	}
	sc, err := seccomp.Seccomp()
	if err != nil {
		errs = append(errs, err)
	}
	ns, err := namespace.Namespace()
	if err != nil {
		errs = append(errs, err)
	}
	cg, err := cgroups.Cgroups()
	if err != nil {
		errs = append(errs, err)
	}
	fs, err := storagedriver.Filesystem()
	if err != nil {
		errs = append(errs, err)
	}
	lsf := container.LinuxSecurityFeature{
		Credential:   container.Credential{},
		Capabilities: caps,
		LSM: container.LSM{
			Apparmor: aa,
			SELinux:  se,
		},
		Seccomp:    sc,
		Namespace:  ns,
		CGroups:    cg,
		Filesystem: fs,
	}
	return lsf, errors.Join(errs...)
}

func Auto() (container.Env, error) {
	var errs []error
	basic, err := Basic()
	if err != nil {
		errs = append(errs, err)
	}
	k, err := kernel.Kernel()
	if err != nil {
		errs = append(errs, err)
	}
	lsf, err := LinuxSecurityFeature()
	if err != nil {
		errs = append(errs, err)
	}
	env := container.Env{
		Basic:                basic,
		Kernel:               k,
		LinuxSecurityFeature: lsf,
		Cluster:              container.Cluster{}, //TODO
		Advance:              container.Advance{}, //TODO
	}
	return env, errors.Join(errs...)
}
