package vul

import (
	"github.com/ctrsploit/ctrsploit/prerequisite"
	"github.com/ctrsploit/ctrsploit/prerequisite/namespace"
)

type networkNamespaceHostLevel struct {
	BaseVulnerability
}

var (
	NetworkNamespaceHostLevel = networkNamespaceHostLevel{
		BaseVulnerability{
			Name:        "network_namespace_host_level",
			Description: "The network namespace of the host is shared",
			CheckSecPrerequisites: prerequisite.Prerequisites{
				&namespace.NetworkNamespaceLevelHost,
			},
			ExploitablePrerequisites: nil,
		}}
)

func (v networkNamespaceHostLevel) Exploit() {
	// TODO
}
