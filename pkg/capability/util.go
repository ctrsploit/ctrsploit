package capability

import (
	"fmt"
	"github.com/containerd/containerd/pkg/cap"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
)

func getCapability(pathStatus string, capType cap.Type) (caps uint64, err error) {
	f, err := os.Open(pathStatus)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	defer f.Close()
	capsMap, err := cap.ParseProcPIDStatus(f)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	caps, _ = capsMap[capType]
	return
}

func GetCapabilityByPid(pid string, capType cap.Type) (caps uint64, err error) {
	return getCapability(fmt.Sprintf("/proc/%s/status", pid), capType)
}

func GetPid1Capability(capType cap.Type) (caps uint64, err error) {
	return GetCapabilityByPid("1", capType)
}

func GetCurrentCapability(capType cap.Type) (caps uint64, err error) {
	return GetCapabilityByPid("self", capType)
}
