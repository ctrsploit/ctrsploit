package kernel

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

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
	return parseVersion(string(releaseBytes[:nullIndex]))
}

func parseVersion(s string) (*semver.Version, error) {
	// 1. first attempt to parse directly
	v, err := semver.NewVersion(s)
	if err == nil {
		return v, nil
	}

	// 2. if failed, sanitize and normalize
	sanitized := sanitizeKernelVersion(s)

	// 3. try parsing again
	v, err = semver.NewVersion(sanitized)
	if err != nil {
		return nil, fmt.Errorf("parse kernel version string '%s' (sanitized: '%s'): %w", s, sanitized, err)
	}
	return v, nil
}

// kernelVersionRegex used to extract Major.Minor.Patch and the rest part
// It can match:
// 4.18 -> 4, 18, 0, ""
// 4.18.0 -> 4, 18, 0, ""
// 4.18.0-147... -> 4, 18, 0, "-147..."
// 4.18.0.147... -> 4, 18, 0, ".147..."
var kernelVersionRegex = regexp.MustCompile(`^(\d+)(?:\.(\d+))?(?:\.(\d+))?(.*)$`)

func sanitizeKernelVersion(s string) string {
	// remove leading and trailing whitespace
	s = strings.TrimSpace(s)

	// replace SemVer disallowed underscore '_' with hyphen '-'
	// to fix "invalid semantic version" issue caused by e.g. x86_64
	s = strings.ReplaceAll(s, "_", "-")

	matches := kernelVersionRegex.FindStringSubmatch(s)
	if matches == nil {
		return s // If it cannot be matched at all, return it as is and let the semver library report an error
	}

	major := matches[1]
	minor := matches[2]
	patch := matches[3]
	rest := matches[4]

	// if minor or patch is missing, set them to "0"
	if minor == "" {
		minor = "0"
	}
	if patch == "" {
		patch = "0"
	}

	base := fmt.Sprintf("%s.%s.%s", major, minor, patch)

	if rest != "" {
		// if rest starts with ".", change it to "-" to comply with SemVer
		// if it already starts with "-" or "+", keep it as is
		if strings.HasPrefix(rest, ".") {
			rest = "-" + strings.TrimPrefix(rest, ".")
		} else if !strings.HasPrefix(rest, "-") && !strings.HasPrefix(rest, "+") {
			// if it starts directly with digits or letters (rare case), force add hyphen
			rest = "-" + rest
		}
	}

	return base + rest
}
