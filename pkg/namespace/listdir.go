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

// ReadInodeNumberMapUnderProc return map[path]ino by walk procfs, skip /proc/pid
func ReadInodeNumberMapUnderProc(proc string) (inoMap map[string]int, err error) {
	inoMap = make(map[string]int)
	err = filepath.Walk(proc, func(path string, info fs.FileInfo, err error) error {
		//log.Logger.Debugf("walk %s", path)
		if err != nil {
			//awesome_error.CheckDebug(err)
			return nil
		}
		if path == proc {
			return nil
		}
		if path == filepath.Join(proc, "sys") && info.IsDir() {
			return filepath.SkipDir
		}
		for _, p := range maskOrReadonlyPath {
			if path == filepath.Join(proc, p) {
				if info.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
			}
		}
		if _, err := strconv.Atoi(info.Name()); err == nil && filepath.Join(proc, info.Name()) == path && info.IsDir() {
			return filepath.SkipDir
		}
		//log.Logger.Debugf("stat %s", path)
		ino := int(info.Sys().(*syscall.Stat_t).Ino)
		inoMap[path] = ino
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
	//log.Logger.Debugf("ino map under /proc: %+v", inoMap)
	for _, ino := range inoMap {
		inoList = append(inoList, ino)
	}
	return
}
