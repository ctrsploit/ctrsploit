package where

import (
	"github.com/ctrsploit/ctrsploit/internal/colorful"
	"github.com/ctrsploit/ctrsploit/pkg/where"
	"github.com/ssst0n3/awesome_libs"
)

const CommandWhereName = "where"

const (
	tplInContainer = `
===========Container===========
{.in}  Is in Container
`
	tplInDocker = `
===========Docker===========
{.dockerenv}  .dockerenv exists
{.rootfs}  rootfs contains 'docker'	
{.cgroups}  cgroups contains 'docker'
{.hosts}  the mount source of /etc/hosts contains 'docker'	
{.hostname}  hostname match regex ^[0-9a-f]{12}$
---
{.in}  => Is in docker
`
	tplInK8s = `
===========k8s===========
{.secret}  {.secret_path} exists
{.hostname}  hostname match k8s pattern
{.hosts}  the mount source of /etc/hosts contains 'pods'
{.cgroups}  contains 'kubepods'
---
{.in}  => is in k8s
`
)

func Container() (err error) {
	c := where.Container{}
	in, err := c.IsIn()
	if err != nil {
		return
	}
	info := awesome_libs.Format(tplInContainer, awesome_libs.Dict{
		"in": colorful.TickOrBallot(in),
	})
	print(info)
	return
}

func Docker() (err error) {
	d := where.Docker{}
	in, err := d.IsIn()
	if err != nil {
		return
	}
	info := awesome_libs.Format(tplInDocker, awesome_libs.Dict{
		"dockerenv": colorful.TickOrBallot(d.DockerEnvFileExists),
		"rootfs":    colorful.TickOrBallot(d.RootfsContainsDocker),
		"cgroups":   colorful.TickOrBallot(d.CgroupContainsDocker),
		"hosts":     colorful.TickOrBallot(d.HostsMountSourceContainsDocker),
		"hostname":  colorful.TickOrBallot(d.HostnameMatchPattern),
		"in":        colorful.TickOrBallot(in),
	})
	print(info)
	return
}

func K8s() (err error) {
	k := where.K8s{}
	in, err := k.IsIn()
	if err != nil {
		return
	}
	info := awesome_libs.Format(tplInK8s, awesome_libs.Dict{
		"secret":      colorful.TickOrBallot(k.DirSecretsExists),
		"secret_path": where.PathDirSecretsExists,
		"hostname":    colorful.TickOrBallot(k.HostnameMatchPattern),
		"hosts":       colorful.TickOrBallot(k.HostsMountSourceContainsPods),
		"cgroups":     colorful.TickOrBallot(k.CgroupContainsKubepods),
		"in":          colorful.TickOrBallot(in),
	})
	print(info)
	return
}
