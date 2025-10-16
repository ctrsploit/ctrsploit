package kernel

import (
	"github.com/ctrsploit/ctrsploit/env/kernel/sysctl"
	"github.com/ctrsploit/ctrsploit/pkg/kernel"
	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const CommandName = "kernel"

func Kernel() (k container.Kernel, err error) {
	k.CompiledDate, err = kernel.GetCompiledDate()
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	k.Sysctl, err = sysctl.Sysctl()
	if err != nil {
		return
	}
	return
}
