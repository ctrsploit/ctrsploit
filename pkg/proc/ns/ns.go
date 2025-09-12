package ns

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/ssst0n3/awesome_libs/awesome_error"
)

var (
	mntNsRegexp = regexp.MustCompile(`^mnt:\[(\d+)]$`)
)

// GetInodeNumber returns inode number of /proc/[pid]/ns/[ns]
func GetInodeNumber(path string) (ino uint64, err error) {
	link, err := os.Readlink(path)
	if err != nil {
		// do not print err
		return
	}
	matches := mntNsRegexp.FindStringSubmatch(link)
	if len(matches) != 2 {
		awesome_error.CheckErr(err)
		return
	}
	ino, err = strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		// unexpected error
		err = fmt.Errorf("failed to convert inode string %q to int for: %w", matches[1], err)
		awesome_error.CheckErr(err)
		return
	}
	return
}
