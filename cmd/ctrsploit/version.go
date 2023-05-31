package main

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/version"
	"github.com/urfave/cli/v2"
)

var versionCommand = &cli.Command{
	Name:    "version",
	Aliases: []string{},
	Usage:   "Show the sploit version information",
	Action: func(context *cli.Context) error {
		fmt.Println(version.DefaultVer())
		return nil
	},
}
