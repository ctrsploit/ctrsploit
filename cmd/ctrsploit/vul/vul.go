package vul

import (
	cve_2016_8867 "github.com/ctrsploit/ctrsploit/vul/cve-2016-8867"
	cve_2020_15257 "github.com/ctrsploit/ctrsploit/vul/cve-2020-15257"
	cve_2025_47290 "github.com/ctrsploit/ctrsploit/vul/cve-2025-47290"
	"github.com/urfave/cli/v2"
)

var Command = &cli.Command{
	Name:    "vul",
	Aliases: []string{"v"},
	Usage:   "check security inside a container",
	Subcommands: []*cli.Command{
		cve_2020_15257.VulCmd,
		cve_2025_47290.VulCmd,
		cve_2016_8867.VulCmd,
	},
}
