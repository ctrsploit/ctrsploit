package internal

import (
	"github.com/urfave/cli/v2"
)

func Command2App(command *cli.Command) (app *cli.App) {
	return &cli.App{
		Name:     command.Name,
		Usage:    command.Usage,
		Action:   command.Action,
		Commands: command.Subcommands,
		Flags:    command.Flags,
		Before:   command.Before,
	}
}
