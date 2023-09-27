package env

import (
	"github.com/ctrsploit/ctrsploit/env"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:    env.SubCommandName,
	Aliases: []string{"e"},
	Usage:   "gather information",
	Subcommands: []*cli.Command{
		Auto,
		WhereCommand,
		graphdriverCommand,
		CgroupsCommand,
		capabilityCommand,
		SeccompCommand,
		ApparmorCommand,
		SelinuxCommand,
		fdiskCommand,
		kernelCommand,
		NamespaceCommand,
	},
}
