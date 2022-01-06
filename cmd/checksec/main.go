package main

import (
	"github.com/ctrsploit/ctrsploit/cmd/ctrsploit/checksec"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:     name,
		Usage:    usage,
		Action:   checksec.Command.Action,
		Commands: checksec.Command.Subcommands,
		Flags:    checksec.Command.Flags,
		Before:   checksec.Command.Before,
	}
	err := app.Run(os.Args)
	if err != nil {
		awesome_error.CheckFatal(err)
	}
}
