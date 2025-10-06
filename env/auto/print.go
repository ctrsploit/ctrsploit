package auto

import (
	"fmt"

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
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
)

type Result struct {
	Where      where.Result
	Mountinfo  mountinfo.Result
	Apparmor   apparmor.Result
	SELinux    selinux.Result
	Capability capability.Caps
	Cgroups    cgroups.Result
	Filesystem storagedriver.Result
	Namespace  namespace.Result
	Seccomp    seccomp.Result
}

func Human(machine container.Env) (human Result) {
	human = Result{
		Where:      where.Human(machine.Where),
		Mountinfo:  mountinfo.Human(machine.MountInfo),
		Apparmor:   apparmor.Human(machine.Apparmor),
		SELinux:    selinux.Human(machine.SELinux),
		Capability: capability.Human(machine.Capabilities),
		Cgroups:    cgroups.Human(machine.CGroups),
		Filesystem: storagedriver.Human(machine.Filesystem),
		Namespace:  namespace.Human(machine.Namespace, ""),
		Seccomp:    seccomp.Human(machine.Seccomp),
	}
	return
}

func Print() (err error) {
	machine, err := Auto()
	if err != nil {
		return
	}
	u := result.Union{
		Machine: machine,
		Human:   Human(machine),
	}
	fmt.Println(printer.Printer.Print(u))
	return
}
