package kernel

import (
	"github.com/ctrsploit/ctrsploit/env/kernel/rlimit"
	"github.com/ctrsploit/ctrsploit/env/kernel/sysctl"
	"github.com/ctrsploit/ctrsploit/pkg/kernel"
	spec "github.com/ctrsploit/sploit-spec/pkg/env/container/kernel"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const CommandName = "kernel"

func Kernel() (k spec.Kernel, err error) {
	k.CompiledDate, err = kernel.GetCompiledDate()
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	k.Sysctl, _ = sysctl.Sysctl()
	k.Rlimit, err = rlimit.Rlimit()
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	return
}
