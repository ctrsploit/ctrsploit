package cpusec

import (
	"bufio"
	"os"
	"strings"
)

// parseCPUInfo reads /proc/cpuinfo and returns the flag tokens from the first
// CPU entry. On x86-64 this is the colon-separated value of the "flags" line;
// on arm64 it is the "Features" line. Flags are uniform across CPUs online at
// boot, so the first entry is authoritative. Returns an empty slice (no error)
// if no flags/Features line is found — callers treat that as "no cpuinfo
// signal" and fall back to the kernel config.
func parseCPUInfo() (flags []string, err error) {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	// /proc/cpuinfo lines are short; bump the buffer only if needed.
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		ln := sc.Text()
		var key, val string
		// Lines look like "flags\t\t: fpu vme ... smep smap ..." or
		// "Features\t: fp asimd ... paca pacg ...". Split on the first colon.
		if i := strings.IndexByte(ln, ':'); i >= 0 {
			key = strings.TrimSpace(ln[:i])
			val = strings.TrimSpace(ln[i+1:])
		} else {
			continue
		}
		if key == "flags" || key == "Features" {
			return strings.Fields(val), nil
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}
