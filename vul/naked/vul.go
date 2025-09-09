package naked

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/apparmor"
	"github.com/ctrsploit/ctrsploit/prerequisite/seccomp"
	"github.com/ctrsploit/ctrsploit/prerequisite/selinux"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
)

var (
	CheckSecCmd = app.Vul2ChecksecCmd(&Vul, nil, nil)
	VulCmd      = app.Vul2VulCmd(&Vul, nil, nil, nil, true)
)

type vulnerability struct {
	vul.BaseVulnerability
}

var (
	Vul = vulnerability{
		BaseVulnerability: vul.BaseVulnerability{
			Name: "naked",
			Description: "We call containers running without seccomp, AppArmor, or SELinux enabled 'naked containers', " +
				"which leaves them highly vulnerable to kernel exploits and potential container escapes",
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites: prerequisite.And(
				&seccomp.Disabled,
				&selinux.Disabled,
				&apparmor.Disabled,
			),
			ExploitablePrerequisites: nil,
		},
	}
)
