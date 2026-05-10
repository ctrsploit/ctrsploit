package pipeprimitive

import (
	"fmt"
	"os"
	"syscall"
)

const dropCachesAll = "3\n"

var (
	dropCachesPath = "/proc/sys/vm/drop_caches"
	syncFS         = syscall.Sync
)

func Clean(syncBeforeDrop bool) error {
	if syncBeforeDrop {
		syncFS()
	}
	if err := os.WriteFile(dropCachesPath, []byte(dropCachesAll), 0); err != nil {
		return fmt.Errorf("drop page cache by writing %q to %s: %w", dropCachesAll, dropCachesPath, err)
	}
	return nil
}
