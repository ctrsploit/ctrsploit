package vul

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/ctrsploit/ctrsploit/prerequisite"
)

type Vulnerability interface {
	Info()
	// CheckSec whether vulnerability exists
	CheckSec() (bool, error)
	// Exploitable whether vulnerability can be exploited,
	// will be called automatically before Exploit()
	Exploitable(vulnerabilityExists bool) (bool, error)
	Exploit()
}

type BaseVulnerability struct {
	Description              string
	CheckSecPrerequisites    prerequisite.Prerequisites
	ExploitablePrerequisites prerequisite.Prerequisites
}

func (v BaseVulnerability) Info() {
	log.Logger.Info(v.Description)
}

func (v BaseVulnerability) CheckSec() (vulnerabilityExists bool, err error) {
	vulnerabilityExists, err = v.CheckSecPrerequisites.Satisfied()
	if err != nil {
		return
	}
	return
}

func (v BaseVulnerability) Exploitable(vulnerabilityExists bool) (satisfied bool, err error) {
	prerequisiteVulnerabilityExists := prerequisite.VulnerabilityExists(vulnerabilityExists)
	v.ExploitablePrerequisites = append([]prerequisite.Interface{prerequisiteVulnerabilityExists}, v.ExploitablePrerequisites...)
	satisfied, err = v.ExploitablePrerequisites.Satisfied()
	if err != nil {
		return
	}
	return
}
