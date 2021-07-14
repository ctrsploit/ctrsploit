package seccomp

import (
	"bufio"
	"os"
	"strings"
)

// parseStatusFile
// borrow from https://github.com/opencontainers/runc/blob/v1.0.0-rc91/libcontainer/seccomp/seccomp_linux.go#L243
func parseStatusFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	status := make(map[string]string)

	for s.Scan() {
		text := s.Text()
		parts := strings.Split(text, ":")

		if len(parts) <= 1 {
			continue
		}

		status[parts[0]] = parts[1]
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return status, nil
}
