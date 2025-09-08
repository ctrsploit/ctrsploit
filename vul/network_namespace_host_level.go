package vul

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/namespace"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v2"
)

type networkNamespaceHostLevel struct {
	vul.BaseVulnerability
}

var (
	NetworkNamespaceHostLevel = networkNamespaceHostLevel{
		vul.BaseVulnerability{
			Name:        "host_net_ns",
			Description: "The network namespace of the host is shared",
			ExeEnv: exeenv.ExeEnv{
				Env:     exeenv.InContainer,
				Check:   exeenv.InContainer,
				Exploit: exeenv.InContainer,
			},
			CheckSecPrerequisites:    &namespace.NetworkNamespaceLevelHost,
			ExploitablePrerequisites: nil,
		}}
)

func (v networkNamespaceHostLevel) Exploit(context *cli.Context) (err error) {
	// TODO
	return
}
