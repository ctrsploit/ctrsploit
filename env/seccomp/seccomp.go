package seccomp

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ctrsploit/ctrsploit/pkg/seccomp"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const (
	CommandName = "seccomp"
)

type Result struct {
	Name      result.Title
	Kernel    item.Bool  `json:"kernel"`
	Container item.Bool  `json:"container"`
	Mode      item.Short `json:"mode"`
}

func (r Result) String() (s string) {
	s += internal.Print(r.Name, r.Kernel, r.Container)
	if r.Container.Result {
		s += internal.Print(r.Mode)
	}
	return
}

// Seccomp
// reference: https://lwn.net/Articles/656307/
func Seccomp() (err error) {
	seccompMode, _, err := seccomp.GetStatus()
	awesome_error.CheckFatal(err)
	r := Result{
		Name: result.Title{
			Name: "Seccomp",
		},
		Kernel: item.Bool{
			Name:        "Kernel Supported",
			Description: "",
			Result:      seccomp.CheckSupported(),
		},
		Container: item.Bool{
			Name:        "Container Enabled",
			Description: "",
			Result:      seccompMode > 0,
		},
		Mode: item.Short{
			Name:        "Mode",
			Description: "",
			Result:      seccompMode.String(),
		},
	}
	fmt.Println(r)
	return
}
