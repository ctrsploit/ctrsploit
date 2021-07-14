package tmp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

func TestRel(t *testing.T) {
	path, err := filepath.Rel(filepath.Join("/root", "/proc"), "/root/proc")
	assert.NoError(t, err)
	assert.Equal(t, ".", path)
}

func checkMountDestination(rootfs, dest string) error {
	invalidDestinations := []string{
		"/proc",
	}
	// White list, it should be sub directories of invalid destinations
	validDestinations := []string{
		// These entries can be bind mounted by files emulated by fuse,
		// so commands like top, free displays stats in container.
		"/proc/cpuinfo",
		"/proc/diskstats",
		"/proc/meminfo",
		"/proc/stat",
		"/proc/swaps",
		"/proc/uptime",
		"/proc/loadavg",
		"/proc/net/dev",
	}
	for _, valid := range validDestinations {
		path, err := filepath.Rel(filepath.Join(rootfs, valid), dest)
		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}
	}
	for _, invalid := range invalidDestinations {
		path, err := filepath.Rel(filepath.Join(rootfs, invalid), dest)
		if err != nil {
			return err
		}
		if path != "." && !strings.HasPrefix(path, "..") {
			return fmt.Errorf("%q cannot be mounted because it is located inside %q", dest, invalid)
		}
	}
	return nil
}

func TestCheckMountDestOnProcChroot(t *testing.T) {
	dest := "/rootfs/proc/"
	err := checkMountDestination("/rootfs", dest)
	if err == nil {
		t.Fatal("destination inside proc when using chroot should not return an error")
	}
}
