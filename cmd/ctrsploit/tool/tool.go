package tool

import (
	"github.com/ctrsploit/ctrsploit/tool"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:     tool.SubCommandName,
	Aliases:  []string{"t"},
	Usage:    "convenience tools for penetration testing",
	Commands: []*cli.Command{
		// add small tools here
	},
}
