package kernel

import (
	"bytes"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"golang.org/x/sys/unix"
)

func Version() (*semver.Version, error) {
	var uts unix.Utsname
	if err := unix.Uname(&uts); err != nil {
		return nil, fmt.Errorf("get uts name: %w", err)
	}
	releaseBytes := uts.Release[:]
	nullIndex := bytes.IndexByte(releaseBytes, 0)
	if nullIndex == -1 {
		nullIndex = len(releaseBytes)
	}
	kernelVersionStr := string(releaseBytes[:nullIndex])
	version, err := semver.NewVersion(kernelVersionStr)
	if err != nil {
		return nil, fmt.Errorf("parse kernel version string '%s': %w", kernelVersionStr, err)
	}
	return version, nil
}
