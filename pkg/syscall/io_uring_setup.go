package syscall

import (
	"github.com/ctrsploit/ctrsploit/pkg/version/docker"
	"github.com/ctrsploit/ctrsploit/pkg/version/kernel"
	"github.com/ctrsploit/ctrsploit/pkg/version/runc"
	"github.com/ctrsploit/ctrsploit/pkg/version/version"
	"golang.org/x/sys/unix"
)

var (
	// IOURingSetup
	// kernel enable:   https://elixir.bootlin.com/linux/v5.1-rc1/source/include/linux/syscalls.h#L318
	// seccomp disable: https://github.com/moby/moby/pull/46762
	// seccomp enable:  http://github.com/moby/moby/pull/39415
	IOURingSetup = Syscall{
		Number: unix.SYS_IO_URING_SETUP,
		Version: map[version.TypeState][]map[version.TypeSoftware][]version.Version{
			version.StateDisable: {
				// docker disables syscall
				{
					// all
					version.SoftwareKernel: kernel.Versions.Values(),
					version.SoftwareDocker: append(
						// <= 19.03.15
						docker.BeforeWhitelistIoUring,
						// >= 25.0.0-beta.1
						docker.CommitBlockIoUring...,
					),
					// all
					version.SoftwareRunc: runc.Versions,
				},
				// kernel does not support syscall, and
				// runc(static) does not support ENOSYS
				{
					// < 5.0
					version.SoftwareKernel: kernel.BeforeIntroduceIoUring,
					// all
					version.SoftwareDocker: docker.Versions.Values(),
					// static, and <= 1.0.0-rc92
					version.SoftwareRunc: runc.StaticBeforeSupportEnosys,
				},
			},
			version.StateValid:       {},
			version.StateUnsupported: {},
		},
	}
)
