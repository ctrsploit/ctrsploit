package env

import (
	"github.com/ctrsploit/ctrsploit/env/version"
	"github.com/urfave/cli/v2"
)

var DockerVersion = &cli.Command{
	Name:    "docker-version",
	Aliases: []string{"dv"},
	Usage:   "guess dockerd version range",
	Action: func(context *cli.Context) (err error) {
		version.Docker()
		return
	},
}
