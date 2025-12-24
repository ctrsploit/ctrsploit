package proc_root

import (
	"os"

	"github.com/ctrsploit/ctrsploit/prerequisite/namespace"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v3"
)

var (
	aliases     = []string{"proc"}
	CheckSecCmd = app.Vul2ChecksecCmd(&Vul, aliases, nil)
	ExploitCmd  = app.Vul2ExploitCmd(&Vul, aliases, nil, true)
	VulCmd      = app.Vul2VulCmd(&Vul, aliases, nil, nil, true)
)

var (
	Vul = vulnerability{
		BaseVulnerability: vul.BaseVulnerability{
			Name:        "host-pid-proc-root",
			Description: "escape by abusing host pid ns via /proc/[pid]/root",
			Level:       vul.LevelHigh,
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites:    &namespace.PidNamespaceLevelHost,
			ExploitablePrerequisites: nil,
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

func (v *vulnerability) Exploit(cmd *cli.Command) (err error) {
	err = v.BaseVulnerability.Exploit(cmd)
	if err != nil {
		return
	}
	return Exploit(os.Stdin, os.Stdout, os.Stderr)
}
