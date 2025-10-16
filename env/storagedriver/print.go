package storagedriver

import (
	"fmt"

	"github.com/ctrsploit/sploit-spec/pkg/env/container"
	"github.com/ctrsploit/sploit-spec/pkg/printer"
	"github.com/ctrsploit/sploit-spec/pkg/result"
	"github.com/ctrsploit/sploit-spec/pkg/result/item"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

type ResultItem struct {
	Name    result.Title `json:"name"`
	Type    item.Short   `json:"type"`
	Enabled item.Bool    `json:"enabled"`
	Used    item.Bool    `json:"used"`
	Number  item.Long    `json:"number"`
	Rootfs  item.Long    `json:"rootfs"`
}

type Result map[string]ResultItem

func Human(machine container.Filesystem) (human Result) {
	human = Result{
		"storage-driver": ResultItem{
			Name: result.Title{
				Name: "Storage Driver",
			},
			Type: item.Short{
				Name:        "Type",
				Description: "",
				Result:      machine.StorageDriver.Type.String(),
			},
			Enabled: item.Bool{
				Name:        "Enabled",
				Description: "",
				Result:      machine.StorageDriver.Enabled,
			},
			Used: item.Bool{
				Name:        "Used",
				Description: "",
				Result:      machine.StorageDriver.Used,
			},
			Number: item.Long{
				Name:        "The number of graph driver mounted",
				Description: "equal to the number of containers",
				Result:      fmt.Sprintf("%d", machine.StorageDriver.Number),
			},
			Rootfs: item.Long{
				Name:        "The host path of container's rootfs",
				Description: "",
				Result:      machine.StorageDriver.Rootfs,
			},
		},
	}
	return
}

func Print() (err error) {
	machine, err := Filesystem()
	u := result.Union{
		Machine: machine,
		Human:   Human(machine),
	}
	fmt.Println(printer.Printer.Print(u))
	// handle error after printing, to avoid missing output
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	return
}
