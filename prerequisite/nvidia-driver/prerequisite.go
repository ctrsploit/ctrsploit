package nvidia_driver

import (
	"fmt"
	"os"

	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
	"golang.org/x/sys/unix"
)

type deviceNode struct {
	prerequisite.BasePrerequisite
	path string
}

func newDeviceNode(path string) *deviceNode {
	return &deviceNode{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name:   fmt.Sprintf("NVIDIA device node %s", path),
			Info:   fmt.Sprintf("%s exists, is a character device, and is readable+writable via access(2); the node is not opened and no ioctl is issued", path),
			ExeEnv: exeenv.InContainer | exeenv.InHost,
		},
		path: path,
	}
}

func (p *deviceNode) Check() (bool, error) {
	return p.CheckTemplate(func() {
		fi, err := os.Stat(p.path)
		if err != nil {
			log.Logger.Debugf("nvidia-driver: cannot stat %s: %v", p.path, err)
			return
		}
		if fi.Mode()&os.ModeCharDevice == 0 {
			log.Logger.Debugf("nvidia-driver: %s is not a character device", p.path)
			return
		}
		if err := unix.Access(p.path, unix.R_OK|unix.W_OK); err != nil {
			log.Logger.Debugf("nvidia-driver: %s is not r/w accessible: %v", p.path, err)
			return
		}
		p.Satisfied = true
	})
}

func newDeviceNodes(paths ...string) prerequisite.Set {
	sets := make([]prerequisite.Set, 0, len(paths))
	for _, path := range paths {
		sets = append(sets, newDeviceNode(path))
	}
	return prerequisite.And(sets...)
}

//goland:noinspection GoSnakeCaseUsage
var (
	VersionVulnerable = newVersion()

	DeviceNodes = newDeviceNodes(
		"/dev/nvidiactl",
		"/dev/nvidia0",
		"/dev/nvidia-uvm",
	)

	Vulnerable = prerequisite.And(VersionVulnerable, DeviceNodes)
)
