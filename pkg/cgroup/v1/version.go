package v1

import (
	"github.com/ctrsploit/ctrsploit/pkg/cgroup/version"
)

func (c CgroupV1) GetVersion() version.Version {
	return version.V1
}
