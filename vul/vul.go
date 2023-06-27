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
	Output()
}

type BaseVulnerability struct {
	Name                     string                     `json:"name"`
	Description              string                     `json:"description"`
	VulnerabilityExists      bool                       `json:"vulnerabilityExists"`
	CheckSecHaveRan          bool                       `json:"-"`
	CheckSecPrerequisites    prerequisite.Prerequisites `json:"-"`
	ExploitablePrerequisites prerequisite.Prerequisites `json:"-"`
}

func (v BaseVulnerability) Info() {
	log.Logger.Info(v.Description)
}

func (v BaseVulnerability) CheckSec() (vulnerabilityExists bool, err error) {
	vulnerabilityExists, err = v.CheckSecPrerequisites.Satisfied()
	if err != nil {
		return
	}
	v.VulnerabilityExists = vulnerabilityExists
	v.CheckSecHaveRan = true
	return
}

func (v BaseVulnerability) Output() {

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
