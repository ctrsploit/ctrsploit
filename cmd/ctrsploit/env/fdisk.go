package env

import (
	"github.com/ctrsploit/ctrsploit/env"
	"github.com/urfave/cli/v2"
)

var Fdisk = &cli.Command{
	Name:    env.CommandFdiskName,
	Aliases: []string{"f"},
	Usage:   "like linux command fdisk or lsblk // TODO",
	Action: func(context *cli.Context) (err error) {
		err = env.Fdisk()
		if err != nil {
			return
		}
		return
	},
}
