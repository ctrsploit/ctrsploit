package env

import (
	"github.com/ctrsploit/ctrsploit/env/selinux"
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/urfave/cli/v2"
)

const (
	CommandNameSelinux = "selinux"
)

var Selinux = &cli.Command{
	Name:    CommandNameSelinux,
	Aliases: []string{"se"},
	Usage:   "show the selinux info",
	Action: func(context *cli.Context) (err error) {
		log.Logger.Debug("")
		err = selinux.Selinux()
		if err != nil {
			return
		}
		return
	},
}
