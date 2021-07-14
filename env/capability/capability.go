package capability

import (
	"ctrsploit/log"
	"ctrsploit/pkg/capability"
	"fmt"
	"github.com/containerd/containerd/pkg/cap"
)

const CommandName = "capability"

func GetPid1Capability() (err error) {
	info := "===========Capability========="
	{	// for pid 1
		caps, err := capability.GetPid1Capability()
		if err != nil {
			return err
		}
		capsParsed, _ := cap.FromBitmap(caps)
		info += fmt.Sprintf("\ncapability of pid 1: %x, %v", caps, capsParsed)
	}
	info += "\n"
	{
		// for current process
		caps, err := capability.GetCurrentCapability()
		if err != nil {
			return err
		}
		capsParsed, _ := cap.FromBitmap(caps)
		info += fmt.Sprintf("\ncapability of current process: %x, %v", caps, capsParsed)
	}
	log.Logger.Info(info)
	return
}
