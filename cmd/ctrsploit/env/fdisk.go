package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env"
	"github.com/urfave/cli/v3"
)

var Fdisk = &cli.Command{
	Name:    env.CommandFdiskName,
	Aliases: []string{"f"},
	Usage:   "like linux command fdisk or lsblk // TODO",
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		err = env.Fdisk()
		if err != nil {
			return
		}
		return
	},
}
