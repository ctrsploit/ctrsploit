package suid

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

const CommandName = "suid"

var defaultPaths = []string{
	"/bin",
	"/sbin",
	"/usr/bin",
	"/usr/sbin",
	"/usr/lib",
	"/usr/libexec",
	"/usr/local/bin",
	"/usr/local/sbin",
}

var defaultSkipDirs = []string{
	"/dev",
	"/proc",
	"/run",
	"/sys",
	"/tmp",
	"/var/run",
	"/var/tmp",
}

type Options struct {
	Paths    []string
	SkipDirs []string
}

type File struct {
	Path       string `json:"path"`
	Mode       string `json:"mode"`
	UID        uint32 `json:"uid"`
	GID        uint32 `json:"gid"`
	Size       int64  `json:"size"`
	Executable bool   `json:"executable"`
	RootOwned  bool   `json:"root_owned"`
}

func DefaultOptions() Options {
	return Options{
		Paths:    append([]string{}, defaultPaths...),
		SkipDirs: append([]string{}, defaultSkipDirs...),
	}
}

func ParsePaths(value string) []string {
	var paths []string
	for _, part := range strings.Split(value, ",") {
		path := strings.TrimSpace(part)
		if path != "" {
			paths = append(paths, path)
		}
	}
	return paths
}

func Find(opts Options) ([]File, error) {
	if len(opts.Paths) == 0 {
		opts.Paths = append([]string{}, defaultPaths...)
	}
	if opts.SkipDirs == nil {
		opts.SkipDirs = append([]string{}, defaultSkipDirs...)
	}

	skipDirs := make(map[string]struct{}, len(opts.SkipDirs))
	for _, dir := range opts.SkipDirs {
		abs, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		skipDirs[filepath.Clean(abs)] = struct{}{}
	}

	seen := make(map[string]struct{})
	var files []File
	var errs []error
	for _, root := range opts.Paths {
		if root == "" {
			continue
		}
		if err := walk(root, skipDirs, seen, &files); err != nil {
			errs = append(errs, err)
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})
	return files, errors.Join(errs...)
}

func walk(root string, skipDirs map[string]struct{}, seen map[string]struct{}, files *[]File) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve %s: %w", root, err)
	}
	absRoot = filepath.Clean(absRoot)

	if _, err := os.Lstat(absRoot); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("stat %s: %w", absRoot, err)
	}

	return filepath.WalkDir(absRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if entry != nil && entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		cleanPath := filepath.Clean(path)
		if entry.IsDir() {
			if cleanPath != absRoot {
				if _, ok := skipDirs[cleanPath]; ok {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return nil
		}
		if !info.Mode().IsRegular() || info.Mode()&os.ModeSetuid == 0 {
			return nil
		}
		if _, ok := seen[cleanPath]; ok {
			return nil
		}
		seen[cleanPath] = struct{}{}

		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return nil
		}
		*files = append(*files, File{
			Path:       cleanPath,
			Mode:       formatMode(info.Mode()),
			UID:        stat.Uid,
			GID:        stat.Gid,
			Size:       info.Size(),
			Executable: info.Mode()&0o111 != 0,
			RootOwned:  stat.Uid == 0,
		})
		return nil
	})
}

func formatMode(mode os.FileMode) string {
	perm := []byte("----------")
	if mode.IsDir() {
		perm[0] = 'd'
	} else if mode&os.ModeSymlink != 0 {
		perm[0] = 'l'
	}
	if mode&0o400 != 0 {
		perm[1] = 'r'
	}
	if mode&0o200 != 0 {
		perm[2] = 'w'
	}
	if mode&0o100 != 0 {
		perm[3] = 'x'
	}
	if mode&0o040 != 0 {
		perm[4] = 'r'
	}
	if mode&0o020 != 0 {
		perm[5] = 'w'
	}
	if mode&0o010 != 0 {
		perm[6] = 'x'
	}
	if mode&0o004 != 0 {
		perm[7] = 'r'
	}
	if mode&0o002 != 0 {
		perm[8] = 'w'
	}
	if mode&0o001 != 0 {
		perm[9] = 'x'
	}
	if mode&os.ModeSetuid != 0 {
		if mode&0o100 != 0 {
			perm[3] = 's'
		} else {
			perm[3] = 'S'
		}
	}
	if mode&os.ModeSetgid != 0 {
		if mode&0o010 != 0 {
			perm[6] = 's'
		} else {
			perm[6] = 'S'
		}
	}
	if mode&os.ModeSticky != 0 {
		if mode&0o001 != 0 {
			perm[9] = 't'
		} else {
			perm[9] = 'T'
		}
	}
	return string(perm)
}
