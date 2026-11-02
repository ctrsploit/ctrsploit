package module

import (
	nvidia_dma_pagecache_write "github.com/ctrsploit/ctrsploit/vul/nvidia-dma-pagecache-write"
	"github.com/urfave/cli/v3"
)

var NvidiaDriver = &cli.Command{
	Name:      "nvidia-driver",
	Aliases:   []string{"nvdriver"},
	Usage:     "nvidia driver related vulnerabilities",
	UsageText: "ctrsploit module nvidia-driver [vul-name]",
	Description: `NVIDIA driver related vulnerabilities, grouped as a logical
module entrypoint. Container Toolkit vulnerabilities stay in the separate
nvidia-container-toolkit group.`,
	Commands: []*cli.Command{
		nvidia_dma_pagecache_write.VulCmd,
	},
}
