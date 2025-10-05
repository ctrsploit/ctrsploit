package cgroups

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_ContainsDocker(t *testing.T) {
	tests := map[string]struct {
		satisfied bool
	}{
		"docker-v28.3.2": {
			// Cgroup Driver: systemd
			// Cgroup Version: 2
			// /proc/1/cgroup
			// 0::/
			satisfied: false,
		},
		"docker-v19.03.13": {
			// Cgroup Driver: cgroupfs
			// /proc/1/cgroup
			//2:cpu,cpuacct:/docker/2403ccfe427f1dc7b6120dccb9ac58c2f8456950380d587a9f33ea7f00dd303d
			//1:name=systemd:/docker/2403ccfe427f1dc7b6120dccb9ac58c2f8456950380d587a9f33ea7f00dd303d
			//0::/system.slice/containerd.service
			satisfied: true,
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		satisfied, err := ContainsDocker.Check()
		assert.NoError(t, err)
		assert.Equal(t, test.satisfied, satisfied)
	})
}
