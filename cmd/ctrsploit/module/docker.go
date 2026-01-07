package module

import (
	docker_sock "github.com/ctrsploit/ctrsploit/vul/shared-socket/docker-sock"
	"github.com/urfave/cli/v3"
)

var Docker = &cli.Command{
	Name:      "docker",
	Aliases:   []string{"d"},
	Usage:     "docker related vulnerabilities",
	UsageText: "ctrsploit module docker [vul-name]",
	Description: `Docker related vulnerabilities, grouped as a logical
module entrypoint.`,
	Commands: []*cli.Command{
		docker_sock.VulCmd,
	},
}
