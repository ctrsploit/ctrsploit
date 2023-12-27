package version

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/seccomp"
	"github.com/ctrsploit/ctrsploit/pkg/syscall"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

type Result struct {
	Name   result.Title `json:"name"`
	Docker item.Short   `json:"docker"`
}

func Docker() {
	seccompMode, _, err := seccomp.GetStatus()
	if err != nil {
		return
	}
	if seccompMode == 0 {
		return
	}
	state := syscall.IOURingSetup.State()
	log.Logger.Debugf("io_uring_setup %v", state)
	v := fmt.Sprintf("dockerd is in %+v", syscall.IOURingSetup.Range())
	r := item.Long{
		Name:        "dockerd-version",
		Description: "dockerd version range",
		Result:      v,
	}
	fmt.Println(printer.Printer.Print(r))
}
