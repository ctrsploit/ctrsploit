package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/apparmor"
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/urfave/cli/v3"
)

const (
	CommandNameApparmor = "apparmor"
)

var Apparmor = &cli.Command{
	Name:    CommandNameApparmor,
	Aliases: []string{"a"},
	Usage:   "show the apparmor info",
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		log.Logger.Debug("")
		err = apparmor.Print()
		if err != nil {
			return
		}
		return
	},
}
