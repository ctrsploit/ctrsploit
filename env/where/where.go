package where

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/where"
	"github.com/ctrsploit/ctrsploit/util"
)

const CommandWhereName = "where"

func Docker() (err error) {
	d := where.Docker{}
	in, err := d.IsIn()
	if err != nil {
		return
	}
	info := "===========Docker========="

	info += "\n.dockerenv exists: "
	info += util.ColorfulTickOrBallot(d.DockerEnvFileExists)
	info += fmt.Sprintf("\nrootfs contains 'docker': %v", util.ColorfulTickOrBallot(d.RootfsContainsDocker))
	info += fmt.Sprintf("\ncgroup contains 'docker': %v", util.ColorfulTickOrBallot(d.CgroupContainsDocker))
	info += fmt.Sprintf("\nthe mount source of /etc/hosts contains 'docker': %v", util.ColorfulTickOrBallot(d.HostsMountSourceContainsDocker))
	info += fmt.Sprintf("\nhostname match regex ^[0-9a-f]{12}$: %v", util.ColorfulTickOrBallot(d.HostnameMatchPattern))
	info += fmt.Sprintf("\n=> is in docker: %v", util.ColorfulTickOrBallot(in))
	fmt.Printf("%s\n\n", info)
	return
}

func K8s() (err error) {
	k := where.K8s{}
	in, err := k.IsIn()
	if err != nil {
		return
	}
	info := "===========k8s========="

	info += fmt.Sprintf("\n%s exists: ", where.PathDirSecretsExists)
	info += util.ColorfulTickOrBallot(k.DirSecretsExists)
	info += fmt.Sprintf("\nhostname match k8s pattern: %s", util.ColorfulTickOrBallot(k.HostnameMatchPattern))
	info += fmt.Sprintf("\nthe mount source of /etc/hosts contains 'pods': %s", util.ColorfulTickOrBallot(k.HostsMountSourceContainsPods))
	info += fmt.Sprintf("\ncgroup contains 'kubepods': %v", util.ColorfulTickOrBallot(k.CgroupContainsKubepods))
	info += fmt.Sprintf("\n=> is in k8s: %v", util.ColorfulTickOrBallot(in))
	fmt.Printf("%s\n\n", info)
	return
}
