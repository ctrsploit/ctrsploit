package version

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/seccomp"
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
	version := ""
	if seccomp.IOURingSetup.Enabled() {
		version = fmt.Sprintf("dockerd is in [%s, %s]", seccomp.IOURingSetup.DockerMinVersion, seccomp.IOURingSetup.DockerMaxVersion)
	} else {
		version = fmt.Sprintf("dockerd is in [, %s], [%s, ]", seccomp.IOURingSetup.DockerMinVersion, seccomp.IOURingSetup.DockerMaxVersion)
	}
	r := item.Long{
		Name:        "dockerd-version",
		Description: "dockerd version range",
		Result:      version,
	}
	fmt.Println(printer.Printer.Print(r))
}
