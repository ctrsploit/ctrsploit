package module

import (
	cve_2021_25748 "github.com/ctrsploit/ctrsploit/vul/cve-2021-25748"
	"github.com/urfave/cli/v3"
)

var IngressNginx = &cli.Command{
	Name:      "ingress-nginx",
	Aliases:   []string{"ingress"},
	Usage:     "ingress-nginx related vulnerabilities",
	UsageText: "ctrsploit module ingress-nginx [vul-name]",
	Description: `Ingress-Nginx related vulnerabilities, grouped as a logical
module entrypoint.`,
	Commands: []*cli.Command{
		cve_2021_25748.VulCmd,
	},
}
