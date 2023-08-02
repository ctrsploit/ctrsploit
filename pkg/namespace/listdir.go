package namespace

import (
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

var (
	maskOrReadonlyPath = []string{
		"keys",
		"acpi",
		"kcore",
		"timer_list",
		"asound",
	}
)

// ListNamespaceDir return map[namespace]ino by reading /proc/<PID>/ns
func ListNamespaceDir(path string) (namespaceInoMap map[string]int, names []string, err error) {
	namespaceInoMap = make(map[string]int)
	entries, err := os.ReadDir(path)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	for _, entry := range entries {
		var link string
		link, err = os.Readlink(filepath.Join(path, entry.Name()))
		if err != nil {
			awesome_error.CheckErr(err)
			return
		}
		inodeNumber := link[strings.Index(link, "[")+1 : strings.Index(link, "]")]
		var number int
		number, err = strconv.Atoi(inodeNumber)
		namespaceInoMap[entry.Name()] = number
		names = append(names, entry.Name())
	}
	return
}

// ReadInodeNumberMapUnderProc return map[path]ino by walk procfs, pass /proc/pid dir
func ReadInodeNumberMapUnderProc(proc string) (inoMap map[string]int, err error) {
	inoMap = make(map[string]int)
	err = filepath.Walk(proc, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if path == proc {
			return nil
		}
		for _, p := range maskOrReadonlyPath {
			if path == filepath.Join(proc, p) {
				return filepath.SkipDir
			}
		}
		// pass net->self, self->pid ...
		if info.Mode() == fs.ModeSymlink {
			return filepath.SkipDir
		}
		if _, err := strconv.Atoi(info.Name()); filepath.Join(proc, info.Name()) == path && err == nil {
			return filepath.SkipDir
		}
		ino := int(info.Sys().(*syscall.Stat_t).Ino)
		inoMap[path] = ino
		if path == filepath.Join(proc, "sys") {
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		awesome_error.CheckWarning(err)
		return
	}
	return
}

func ReadInodeNumberListUnderProc(proc string) (inoList []int, err error) {
	inoMap, err := ReadInodeNumberMapUnderProc(proc)
	if err != nil {
		return
	}
	// log.Logger.Debugf("ino map under /proc: %+v", inoMap)
	for _, ino := range inoMap {
		inoList = append(inoList, ino)
	}
	return
}
