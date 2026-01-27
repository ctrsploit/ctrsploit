package module

import (
	cve_2022_39253 "github.com/ctrsploit/ctrsploit/vul/cve-2022-39253"
	"github.com/urfave/cli/v3"
)

var Git = &cli.Command{
	Name:      "git",
	Aliases:   []string{"g"},
	Usage:     "git related vulnerabilities",
	UsageText: "ctrsploit module git [vul-name]",
	Description: `Git related vulnerabilities, grouped as a logical
module entrypoint.`,
	Commands: []*cli.Command{
		cve_2022_39253.VulCmd,
	},
}
