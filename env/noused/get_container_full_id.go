package noused

import (
	"fmt"

	proccgroup "github.com/ctrsploit/ctrsploit/pkg/proc/cgroup"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

func GetContainerFullId() (ID string, err error) {
	ID, err = proccgroup.GetContainerID()
	if err != nil {
		awesome_error.CheckErr(fmt.Errorf("get container full id: %w", err))
	}
	return
}
