package module

import (
	cve_2024_23650 "github.com/ctrsploit/ctrsploit/vul/cve-2024-23650"
	"github.com/urfave/cli/v3"
)

var Buildkit = &cli.Command{
	Name:      "buildkit",
	Aliases:   []string{"bk"},
	Usage:     "buildkit related vulnerabilities",
	UsageText: "ctrsploit module buildkit [vul-name]",
	Description: `BuildKit related vulnerabilities, grouped as a logical
module entrypoint.`,
	Commands: []*cli.Command{
		cve_2024_23650.VulCmd,
	},
}
