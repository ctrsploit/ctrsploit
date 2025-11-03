package maskedPath

import (
	"fmt"
	"strings"

	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type Masked struct {
	prerequisite.BasePrerequisite
	path string
}

// Check if the given path is masked (i.e., mounted as /dev/null).
// returns true if masked
// returns false if not masked or path not exists.
func (p *Masked) Check() (bool, error) {
	return p.CheckTemplate(func() {
		exists, _ := internal.CheckPathExists(p.path)
		if !exists {
			return
		}
		info, err := mountinfo.GetMountByMountpoint(p.path)
		if err != nil {
			if strings.Contains(err.Error(), mountinfo.ErrMountPointNotFound) {
				return
			} else {
				p.Err = p.WrapErr(fmt.Errorf("getting mountpoint of: %s, %w", p.path, err))
			}
		}
		if info.Root == "/dev/null" {
			p.Satisfied = true
		}
	})
}
