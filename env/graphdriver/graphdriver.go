package graphdriver

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver"
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver/devicemapper"
	"github.com/ctrsploit/ctrsploit/pkg/graphdriver/overlay"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
)

const CommandName = "graphdriver"

type Result struct {
	Name     result.Title `json:"name"`
	Enabled  item.Bool    `json:"enabled"`
	Used     item.Bool    `json:"used"`
	Number   item.Long    `json:"number"`
	HostPath item.Long    `json:"host_path"`
}

func Overlay() (err error) {
	r, err := graphDriver("Overlay", &overlay.Overlay{})
	if err != nil {
		return
	}
	fmt.Println(printer.Printer.Print(r))
	return
}

func DeviceMapper() (err error) {
	r, err := graphDriver("DeviceMapper", &devicemapper.DeviceMapper{})
	if err != nil {
		return
	}
	fmt.Println(printer.Printer.Print(r))
	return
}

func GraphDrivers() (err error) {
	o, err := graphDriver("Overlay", &overlay.Overlay{})
	if err != nil {
		return
	}
	d, err := graphDriver("DeviceMapper", &devicemapper.DeviceMapper{})
	if err != nil {
		return
	}
	r := map[string]Result{
		"overlay":      o,
		"devicemapper": d,
	}
	fmt.Println(printer.Printer.Print(r))
	return
}

func graphDriver(name string, g graphdriver.Interface) (r Result, err error) {
	err = g.Init()
	if err != nil {
		return
	}

	enabled, err := g.IsEnabled()
	if err != nil {
		return
	}

	r = Result{
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

	return
}
