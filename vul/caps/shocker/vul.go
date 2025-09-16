package shocker

import (
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/urfave/cli/v2"
)

var (
	aliases      = []string{"cap_dac_read_search", "open_by_handle_at"}
	flagsExploit = []cli.Flag{
		&cli.IntFlag{
			Name:        "inode",
			DefaultText: "default is 2, (in ext fs, root's inode is 2)",
			Required:    false,
			Value:       2,
		},
		&cli.StringFlag{
			Name:        "reference",
			Aliases:     []string{"r", "ref", "mountFd"},
			DefaultText: "default is /etc/hosts",
			Required:    false,
			Value:       "/etc/hosts",
		},
	}
	ExploitCmd  = app.Vul2ExploitCmd(&Shocker, aliases, flagsExploit, true)
	CheckSecCmd = app.Vul2ChecksecCmd(&Shocker, aliases, nil)
	VulCmd      = app.Vul2VulCmd(&Shocker, aliases, nil, flagsExploit, true)
)
