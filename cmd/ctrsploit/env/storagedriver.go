package env

import (
	"github.com/ctrsploit/ctrsploit/env/storagedriver"
	"github.com/urfave/cli/v2"
)

var StorageDriver = &cli.Command{
	Name:    storagedriver.CommandName,
	Aliases: []string{"sd"},
	Usage:   "detect storage driver type and extend information",
	Action: func(context *cli.Context) (err error) {
		err = storagedriver.Print()
		if err != nil {
			return
		}
		return
	},
}
