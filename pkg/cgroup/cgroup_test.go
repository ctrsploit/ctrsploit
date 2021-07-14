package cgroup

import (
	"github.com/ssst0n3/awesome_libs/log"
	"testing"
)

func TestIsCgroupV2(t *testing.T) {
	log.Logger.Info(IsCgroupV2())
}
