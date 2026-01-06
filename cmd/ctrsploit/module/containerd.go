package module

import (
	cve_2024_40635 "github.com/ctrsploit/ctrsploit/vul/cve-2024-40635"
	"github.com/urfave/cli/v3"
)

var Containerd = &cli.Command{
	Name:      "containerd",
	Aliases:   []string{"c"},
	Usage:     "containerd related vulnerabilities",
	UsageText: "ctrsploit module containerd [vul-name]",
	Description: `Containerd related vulnerabilities, grouped as a logical
module entrypoint.`,
	Commands: []*cli.Command{
		cve_2024_40635.VulCmd,
	},
}
