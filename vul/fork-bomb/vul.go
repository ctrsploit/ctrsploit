package fork_bomb

import (

	"github.com/ctrsploit/ctrsploit/prerequisite/cgroups/pids"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v3"
)

var (
	flagsExploit = []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "ignore exploitable check",
			Value:   false,
		},
	}
	CheckSecCmd = app.Vul2ChecksecCmd(&Vul, nil, flagsExploit)
	ExploitCmd  = app.Vul2ExploitCmd(&Vul, nil, flagsExploit, true)
	VulCmd      = app.Vul2VulCmd(&Vul, nil, nil, flagsExploit, true)
)

type vulnerability struct {
	vul.BaseVulnerability
}

var (
	Vul = vulnerability{
		BaseVulnerability: vul.BaseVulnerability{
			Name:        "fork-bomb",
			Description: "",
			Level:       0,
			ExeEnv:      exeenv.ExeEnv{},
			CheckSecPrerequisites: prerequisite.And(
				&pids.UnlimitedPidsMax,
				// TODO: do more check
				//// bypass RLIMIT_NPROC
				//prerequisite.Or(
				//	// TODO: not in username space
				//	// The RLIMIT_NPROC limit is not enforced for processes that
				//	// have either the CAP_SYS_ADMIN or the CAP_SYS_RESOURCE
				//	// capability, or run with real user ID 0.
				//	// https://github.com/torvalds/linux/blob/v6.17/kernel/fork.c#L2044-L2045
				//	&capability.CapSysAdminBnd,
				//	&capability.CapSysResourceBnd,
				//	&user.EUid0, // may set up ruid
				//	&user.RUid0,
				//),
			),
			ExploitablePrerequisites: nil,
		},
	}
)

func (v *vulnerability) Exploit(cmd *cli.Command) (err error) {
	err = v.BaseVulnerability.Exploit(cmd)
	if err != nil {
		return
	}
	return Exploit()
}
