package graphdriver

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver"
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver/devicemapper"
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver/overlay"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

const CommandName = "graphdriver"

type Result struct {
	Name     result.Title
	Enabled  item.Bool `json:"enabled"`
	Used     item.Bool `json:"used"`
	Number   item.Long `json:"number"`
	HostPath item.Long `json:"host_path"`
}

func (r Result) String() (s string) {
	s += internal.Print(r.Name, r.Enabled)
	if r.Enabled.Result {
		s += internal.Print(r.Used)
		if r.Used.Result {
			s += internal.Print(r.Number, r.HostPath)
		}
	}
	return
}

func Overlay() (err error) {
	return graphDriver("Overlay", &overlay.Overlay{})
}

func DeviceMapper() (err error) {
	return graphDriver("DeviceMapper", &devicemapper.DeviceMapper{})
}

func graphDriver(name string, g graphdriver.Interface) (err error) {
	err = g.Init()
	if err != nil {
		return
	}

	enabled, err := g.IsEnabled()
	if err != nil {
		return
	}

	r := Result{
		Name: result.Title{
			Name: name,
		},
		Enabled: item.Bool{
			Name:        "Enabled",
			Description: "",
			Result:      enabled,
		},
	}
	var used bool
	var number int
	var hostPath string
	if r.Enabled.Result {
		used, err = g.IsUsed()
		if err != nil {
			return
		}
		r.Used = item.Bool{
			Name:        "Used",
			Description: "",
			Result:      used,
		}
		if r.Used.Result {
			number, err = g.Number()
			if err != nil {
				return
			}
			// number = 0 means /sys/module/seccomp/refcnt not exists
			if number > 0 {
				r.Number = item.Long{
					Name:        "The number of graph driver mounted",
					Description: "equal to the number of containers",
					Result:      fmt.Sprintf("%d", number),
				}
			}

			hostPath, err = g.HostPathOfCtrRootfs()
			if err != nil {
				return
			}
			r.HostPath = item.Long{
				Name:        "The host path of container's rootfs",
				Description: "",
				Result:      hostPath,
			}
		}
	}

	fmt.Println(r)
	return
}
