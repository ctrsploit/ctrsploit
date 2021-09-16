package helper

import "github.com/urfave/cli/v2"

var Command = &cli.Command{
	Name:        "helper",
	Aliases:     []string{"he"},
	Usage:       "some helper commands such as local privilege escalation",
	Subcommands: []*cli.Command{},
}
