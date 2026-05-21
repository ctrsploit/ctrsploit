package maskedPath

import (
	"fmt"
	"strings"

	"github.com/ctrsploit/ctrsploit/pkg/fileutil"
	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type Masked struct {
	prerequisite.BasePrerequisite
	Path string
}

// Check if the given Path is masked (i.e., mounted as /dev/null).
// returns true if masked
// returns false if not masked or Path not exists.
func (p *Masked) Check() (bool, error) {
	return p.CheckTemplate(func() {
		exists, _ := fileutil.CheckPathExists(p.Path)
		if !exists {
			return
		}
		info, err := mountinfo.GetMountByMountpoint(p.Path)
		if err != nil {
			if strings.Contains(err.Error(), mountinfo.ErrMountPointNotFound) {
				return
			} else {
				p.Err = p.WrapErr(fmt.Errorf("getting mountpoint of: %s, %w", p.Path, err))
			}
		}
		if info.Root == "/dev/null" {
			p.Satisfied = true
		}
	})
}
