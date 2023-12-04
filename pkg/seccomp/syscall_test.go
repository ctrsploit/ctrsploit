package seccomp

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/version"
	"golang.org/x/sys/unix"
	"testing"
)

func TestVersionRanges_String(t *testing.T) {
	v := VersionRanges{
		{
			Status{Version: version.FirstDockerVersion, Enable: false},
			Status{Version: version.FurtherDockerVersion, Enable: true},
		},
	}
	fmt.Println(v.String())

	fmt.Println(unix.ENOSYS)
}
