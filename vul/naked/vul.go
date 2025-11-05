package naked

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/apparmor"
	"github.com/ctrsploit/ctrsploit/prerequisite/seccomp"
	"github.com/ctrsploit/ctrsploit/prerequisite/selinux"
	"github.com/ctrsploit/ctrsploit/prerequisite/sysctl/user"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

var (
	CheckSecCmd = getCheckSecCmd(Vul.Name, Vul.Description)
	VulCmd      = &cli.Command{
		Name:  Vul.Name,
		Usage: Vul.Description,
		Subcommands: []*cli.Command{
			getCheckSecCmd("checksec", "check vulnerability exists"),
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

var (
	Vul = vulnerability{
		BaseVulnerability: vul.BaseVulnerability{
			Name: "naked",
			Description: "we call containers running without seccomp, AppArmor, or SELinux enabled 'naked containers', " +
				"which leaves them highly vulnerable to kernel exploits and potential container escapes",
			Level: vul.LevelHigh,
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites: prerequisite.And(
				&seccomp.Disabled,
				&selinux.Unconfined,
				&apparmor.Disabled,
				&user.UserNsEnabled,
			),
			ExploitablePrerequisites: nil,
		},
	}
)

func getCheckSecCmd(name, usage string) (cmd *cli.Command) {
	cmd = app.Vul2ChecksecCmd(&Vul, nil, nil)
	cmd.Name = name
	cmd.Usage = usage
	return
}
