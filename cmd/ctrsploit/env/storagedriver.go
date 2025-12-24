package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/storagedriver"
	"github.com/urfave/cli/v3"
)

var StorageDriver = &cli.Command{
	Name:    storagedriver.CommandName,
	Aliases: []string{"sd"},
	Usage:   "detect storage driver type and extend information",
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		err = storagedriver.Print()
		if err != nil {
			return
		}
		return
	},
}
