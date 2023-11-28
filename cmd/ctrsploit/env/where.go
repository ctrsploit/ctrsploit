package env

import (
	"github.com/ctrsploit/ctrsploit/env/where"
	"github.com/urfave/cli/v2"
)

var WhereCommand = &cli.Command{
	Name:    where.CommandName,
	Aliases: []string{"w"},
	Usage:   "detect whether you are in the container, and which type of the container",
	Action: func(context *cli.Context) (err error) {
		err = where.Where()
		if err != nil {
			return
		}
		return
	},
}
