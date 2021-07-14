package capability

import (
	"ctrsploit/log"
	"ctrsploit/pkg/capability"
	"fmt"
	"github.com/containerd/containerd/pkg/cap"
)

const CommandName = "capability"

func GetPid1Capability() (err error) {
	caps, err := capability.GetPid1Capability()
	if err != nil {
		return
	}
	info := "===========Capability========="
	capsParsed, _ := cap.FromBitmap(caps)
	info += fmt.Sprintf("\ncapability of pid 1: %v, %v", caps, capsParsed)
	log.Logger.Info(info)
	return
}
