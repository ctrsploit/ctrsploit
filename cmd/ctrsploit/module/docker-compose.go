package module

import (
	cve_2025_62725 "github.com/ctrsploit/ctrsploit/vul/cve-2025-62725"
	"github.com/urfave/cli/v3"
)

var DockerCompose = &cli.Command{
	Name:      "docker-compose",
	Aliases:   []string{"compose"},
	Usage:     "docker-compose related vulnerabilities",
	UsageText: "ctrsploit module docker-compose [vul-name]",
	Description: `Docker Compose related vulnerabilities, grouped as a logical
module entrypoint.`,
	Commands: []*cli.Command{
		cve_2025_62725.VulCmd,
	},
}
