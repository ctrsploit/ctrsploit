package vul

import (
	"github.com/ctrsploit/ctrsploit/internal/log"
	"github.com/ctrsploit/ctrsploit/prerequisite"
	"github.com/ctrsploit/ctrsploit/prerequisite/vulnerability"
)

type Vulnerability interface {
	// GetName returns a one word name; may be used as command name
	GetName() string
	// GetDescription return usage
	GetDescription() string
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
	VulnerabilityExists      bool                       `json:"vulnerability_exists"`
	CheckSecHaveRan          bool                       `json:"-"`
	CheckSecPrerequisites    prerequisite.Prerequisites `json:"-"`
	ExploitablePrerequisites prerequisite.Prerequisites `json:"-"`
}

func (v *BaseVulnerability) GetName() string {
	return v.Name
}

func (v *BaseVulnerability) GetDescription() string {
	return v.Description
}

func (v *BaseVulnerability) Info() {
	log.Logger.Info(v.Description)
}

func (v *BaseVulnerability) CheckSec() (vulnerabilityExists bool, err error) {
	vulnerabilityExists, err = v.CheckSecPrerequisites.Satisfied()
	if err != nil {
		return
	}
	v.VulnerabilityExists = vulnerabilityExists
	v.CheckSecHaveRan = true
	return
}

func (v *BaseVulnerability) Output() {

}

func (v *BaseVulnerability) Exploitable(vulnerabilityExists bool) (satisfied bool, err error) {
	prerequisiteVulnerabilityExists := vulnerability.Exists(vulnerabilityExists)
	v.ExploitablePrerequisites = append([]prerequisite.Interface{prerequisiteVulnerabilityExists}, v.ExploitablePrerequisites...)
	satisfied, err = v.ExploitablePrerequisites.Satisfied()
	if err != nil {
		return
	}
	return
}
