package module

import (
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "module",
	Aliases:   []string{"m"},
	Usage:     "group vulnerabilities by component or config type",
	UsageText: "ctrsploit module [component|config] [name]",
	Description: `Classify and operate vulnerabilities by logical module
such as kernel, runc, containerd, or config (e.g. capability).`,
	Commands: []*cli.Command{
		Runc,
		Containerd,
		NvidiaContainerToolkit,
		Kernel,
	},
}
