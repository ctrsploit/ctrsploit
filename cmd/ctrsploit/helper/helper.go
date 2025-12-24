package helper

import "github.com/urfave/cli/v3"

var Command = &cli.Command{
	Name:        "helper",
	Aliases:     []string{"he"},
	Usage:       "some helper commands such as local privilege escalation",
	Commands: []*cli.Command{
		CVE_2021_3493_Command,
	},
}
