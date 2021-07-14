package capability

import (
	"github.com/containerd/containerd/pkg/cap"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
)

func getCapability(pathStatus string) (caps uint64, err error) {
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
	caps, _ = capsMap[cap.Effective]
	return
}

func GetPid1Capability() (caps uint64, err error) {
	return getCapability("/proc/1/status")
}

func GetCurrentCapability() (caps uint64, err error) {
	return getCapability("/proc/self/status")
}
