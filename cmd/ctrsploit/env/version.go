package env

import (
	"context"

	"github.com/ctrsploit/ctrsploit/env/version"
	"github.com/urfave/cli/v3"
)

var DockerVersion = &cli.Command{
	Name:    "docker-version",
	Aliases: []string{"dv"},
	Usage:   "guess dockerd version range",
	Action: func(ctx context.Context, cmd *cli.Command) (err error) {
		version.Docker()
		return
	},
}
