package vul

import (
	"github.com/ctrsploit/ctrsploit/prerequisite/namespace"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
)

type networkNamespaceHostLevel struct {
	vul.BaseVulnerability
}

var (
	NetworkNamespaceHostLevel = networkNamespaceHostLevel{
		vul.BaseVulnerability{
			Name:        "host_net_ns",
			Description: "The network namespace of the host is shared",
			CheckSecPrerequisites: prerequisite.Prerequisites{
				&namespace.NetworkNamespaceLevelHost,
			},
			ExploitablePrerequisites: nil,
		}}
)

func (v networkNamespaceHostLevel) Exploit() (err error) {
	// TODO
	return
}
