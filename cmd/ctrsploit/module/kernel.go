package module

import (
	cve_2022_0492 "github.com/ctrsploit/ctrsploit/vul/cve-2022-0492"
	"github.com/urfave/cli/v3"
)

var Kernel = &cli.Command{
	Name:      "kernel",
	Aliases:   []string{"k"},
	Usage:     "kernel related vulnerabilities",
	UsageText: "ctrsploit module kernel [vul-name]",
	Description: `Kernel related vulnerabilities, focused on kernel CVEs
grouped as a logical module entrypoint.`,
	Commands: []*cli.Command{
		cve_2022_0492.VulCmd,
	},
}
