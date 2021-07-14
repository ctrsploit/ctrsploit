package env

import (
	"ctrsploit/env"
	"github.com/urfave/cli/v2"
)

var whereCommand = &cli.Command{
	Name:    env.CommandWhereName,
	Aliases: []string{"w"},
	Usage:   "detect whether you are in the container, and which type of the container",
	Action: func(context *cli.Context) (err error) {
		err = env.Docker()
		if err != nil {
			return
		}
		err = env.K8s()
		if err != nil {
			return
		}
		return
	},
}
