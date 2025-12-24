package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/selinux"
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/urfave/cli/v3"
)

const (
	CommandNameSelinux = "selinux"
)

var Selinux = &cli.Command{
	Name:    CommandNameSelinux,
	Aliases: []string{"se"},
	Usage:   "show the selinux info",
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		log.Logger.Debug("")
		err = selinux.Print()
		if err != nil {
			return
		}
		return
	},
}
