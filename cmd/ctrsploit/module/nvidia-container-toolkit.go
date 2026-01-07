package module

import (
	cve_2024_0132 "github.com/ctrsploit/ctrsploit/vul/cve-2024-0132"
	cve_2025_23266 "github.com/ctrsploit/ctrsploit/vul/cve-2025-23266"
	"github.com/urfave/cli/v3"
)

var NvidiaContainerToolkit = &cli.Command{
	Name:      "nvidia-container-toolkit",
	Aliases:   []string{"nvidia", "nct"},
	Usage:     "nvidia-container-toolkit related vulnerabilities",
	UsageText: "ctrsploit module nvidia-container-toolkit [vul-name]",
	Description: `NVIDIA Container Toolkit related vulnerabilities,
grouped as a logical module entrypoint.`,
	Commands: []*cli.Command{
		cve_2024_0132.VulCmd,
		cve_2025_23266.VulCmd,
	},
}
