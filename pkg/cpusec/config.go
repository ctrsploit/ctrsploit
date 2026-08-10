package cpusec

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"strings"

	"github.com/ctrsploit/ctrsploit/pkg/kernel/uname"
)

// parseKernelConfig reads CONFIG_* lines from the running kernel's config.
// It tries, in order:
//  1. /boot/config-<release>  (release from uname -r)
//  2. /proc/config.gz         (gzip-compressed; available when CONFIG_IKCONFIG+CONFIG_IKCONFIG_PROC)
//
// Returns a map from the CONFIG_* symbol (e.g. "CONFIG_PAGE_TABLE_ISOLATION")
// to its value ("y", "m", "n", or a quoted string), the source path it read,
// and nil error. If neither source exists, it returns an empty map, an empty
// source, and nil error — callers fall back to cpuinfo-only detection.
func parseKernelConfig() (cfg map[string]string, source string, err error) {
	release, relErr := uname.Release()
	if relErr == nil && release != "" {
		path := "/boot/config-" + release
		if c, s, e := readConfigFile(path); e == nil {
			return c, s, nil
		}
	}
	if c, s, e := readConfigGz("/proc/config.gz"); e == nil {
		return c, s, nil
	}
	if relErr != nil {
		return nil, "", fmt.Errorf("uname release: %w", relErr)
	}
	return nil, "", nil
}

// readConfigFile parses a plain-text kernel config at path.
func readConfigFile(path string) (cfg map[string]string, source string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	return parseConfigReader(f, path)
}

// readConfigGz parses a gzip-compressed kernel config at path.
func readConfigGz(path string) (cfg map[string]string, source string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, "", err
	}
	defer gz.Close()
	return parseConfigReader(gz, path)
}

// parseConfigReader consumes a kernel-config stream and returns the CONFIG_*
// map. Lines that are not "CONFIG_FOO=value" (comments, "# CONFIG_FOO is not
// set") are skipped; "is not set" is recorded as value "n" so callers can
// uniformly test cfg["CONFIG_FOO"] == "y".
func parseConfigReader(r interface {
	Read(p []byte) (n int, err error)
}, source string) (cfg map[string]string, src string, err error) {
	cfg = make(map[string]string)
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		ln := sc.Text()
		if strings.HasPrefix(ln, "CONFIG_") {
			if i := strings.IndexByte(ln, '='); i >= 0 {
				cfg[ln[:i]] = strings.TrimSpace(ln[i+1:])
			}
		} else if strings.HasPrefix(ln, "# ") && strings.HasSuffix(ln, " is not set") {
			// "# CONFIG_FOO is not set" → CONFIG_FOO=n
			inner := strings.TrimPrefix(ln, "# ")
			inner = strings.TrimSuffix(inner, " is not set")
			if strings.HasPrefix(inner, "CONFIG_") {
				cfg[inner] = "n"
			}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, "", err
	}
	return cfg, source, nil
}
