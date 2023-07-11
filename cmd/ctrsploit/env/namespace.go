package env

import (
	"github.com/ctrsploit/ctrsploit/env/namespace"
	"github.com/urfave/cli/v2"
)

var NamespaceCommand = &cli.Command{
	Name:    namespace.CommandName,
	Aliases: []string{"n", "ns"},
	Usage:   "check namespace is host ns",
	Action: func(context *cli.Context) (err error) {
		err = namespace.CheckCurrentNamespaceIsHost()
		if err != nil {
			return
		}
		return
	},
}
