package net

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ssst0n3/awesome_libs/awesome_error"
)

/*
/proc/net/unix Lists the UNIX domain sockets present within the system and their status.

The code that generates /proc/net/unix is in the unix_seq_show() function in
[net/unix/af_unix.c](https://elixir.bootlin.com/linux/v6.14.4/source/net/unix/af_unix.c#L3397) in the kernel source.
*/

// ListUnixSocketPath reads /proc/net/unix and returns the paths of the UNIX domain sockets.
func ListUnixSocketPath() (paths []string, err error) {
	file, err := os.Open("/proc/net/unix")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Skip the header line
	if scanner.Scan() {
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			// The last field is the socket path
			if len(fields) > 7 {
				paths = append(paths, fields[len(fields)-1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return paths, nil
}

func FilterUnixSocketByPrefix(prefix string) (paths []string, err error) {
	paths, err = ListUnixSocketPath()
	if err != nil {
		return nil, err
	}

	filteredPathsMap := make(map[string]struct{})
	for _, path := range paths {
		if strings.HasPrefix(path, prefix) {
			filteredPathsMap[path] = struct{}{}
		}
	}

	var filteredPaths []string
	for path := range filteredPathsMap {
		filteredPaths = append(filteredPaths, path)
	}

	return filteredPaths, nil
}

func ContainerdShimAbstractUnixSocketPath(prefix string) (path string, err error) {
	paths, err := FilterUnixSocketByPrefix("@/containerd-shim/")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	if len(paths) == 0 {
		awesome_error.CheckWarning(fmt.Errorf("no containerd-shim abstract unix socket found"))
		return
	}
	path = paths[0]
	return
}
