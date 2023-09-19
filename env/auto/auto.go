package auto

import (
	"github.com/ctrsploit/ctrsploit/env/apparmor"
	"github.com/ctrsploit/ctrsploit/env/capability"
	"github.com/ctrsploit/ctrsploit/env/cgroups"
	"github.com/ctrsploit/ctrsploit/env/graphdriver"
	"github.com/ctrsploit/ctrsploit/env/namespace"
	"github.com/ctrsploit/ctrsploit/env/seccomp"
	"github.com/ctrsploit/ctrsploit/env/selinux"
	"github.com/ctrsploit/ctrsploit/env/where"
)

func Auto() (err error) {
	_ = where.Docker()
	_ = where.K8s()
	_ = apparmor.Apparmor()
	_ = selinux.Selinux()
	_ = capability.Capability()
	_ = cgroups.Version()
	_ = cgroups.ListSubsystems()
	_ = graphdriver.Overlay()
	_ = graphdriver.DeviceMapper()
	_ = namespace.CheckCurrentNamespaceLevel("")
	_ = seccomp.Seccomp()
	return
}
