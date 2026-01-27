package module

import (
	cve_2016_8867 "github.com/ctrsploit/ctrsploit/vul/cve-2016-8867"
	cve_2019_5736 "github.com/ctrsploit/ctrsploit/vul/cve-2019-5736"
	"github.com/urfave/cli/v3"
)

var Runc = &cli.Command{
	Name:      "runc",
	Aliases:   []string{"r"},
	Usage:     "runc related vulnerabilities",
	UsageText: "ctrsploit module runc [vul-name]",
	Description: `Runc related vulnerabilities, grouped as a logical
module entrypoint.`,
	Commands: []*cli.Command{
		cve_2016_8867.VulCmd,
		cve_2019_5736.VulCmd,
	},
}
