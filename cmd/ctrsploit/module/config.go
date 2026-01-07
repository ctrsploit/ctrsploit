package module

import "github.com/urfave/cli/v3"

var Config = &cli.Command{
	Name:      "config",
	Aliases:   []string{"cfg"},
	Usage:     "insecure configuration and misconfiguration issues",
	UsageText: "ctrsploit module config [type] [name]",
	Description: `Entry for insecure configuration (misconfiguration) issues,
such as over-privileged capabilities, weak seccomp profiles, or unsafe
namespace sharing. Concrete subcommands will group these config issues
by type (e.g. capability, seccomp, namespace).`,
	Commands: []*cli.Command{},
}
