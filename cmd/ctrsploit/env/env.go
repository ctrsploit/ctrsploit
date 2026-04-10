package env

import (
	"github.com/ctrsploit/ctrsploit/env"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:    env.SubCommandName,
	Aliases: []string{"e"},
	Usage:   "gather information",
	Commands: []*cli.Command{
		Auto,
		Where,
		Mountinfo,
		StorageDriver,
		Cgroups,
		Capability,
		Seccomp,
		Apparmor,
		Selinux,
		Fdisk,
		Kernel,
		Sysctl,
		Rlimit,
		Namespace,
		DockerVersion,
		Services,
		Upload,
	},
}
