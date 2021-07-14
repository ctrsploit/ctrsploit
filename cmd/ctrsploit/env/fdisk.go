package env

import (
	"ctrsploit/env"
	"github.com/urfave/cli/v2"
)

var fdiskCommand = &cli.Command{
	Name:    env.CommandFdiskName,
	Aliases: []string{"f"},
	Usage:   "like linux command fdisk -l",
	Action: func(context *cli.Context) (err error) {
		env.Fdisk()
		return
	},
}
