package module

import (
	cve_2020_8558 "github.com/ctrsploit/ctrsploit/vul/cve-2020-8558"
	cve_2021_25741 "github.com/ctrsploit/ctrsploit/vul/cve-2021-25741"
	"github.com/urfave/cli/v3"
)

var Kubernetes = &cli.Command{
	Name:      "kubernetes",
	Aliases:   []string{"k8s"},
	Usage:     "kubernetes related vulnerabilities",
	UsageText: "ctrsploit module kubernetes [vul-name]",
	Description: `Kubernetes related vulnerabilities, grouped as a logical
module entrypoint.`,
	Commands: []*cli.Command{
		cve_2020_8558.VulCmd,
		cve_2021_25741.VulCmd,
	},
}
