package status

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Status represents the information in /proc/[pid]/status file.
type Status struct {
	Name       string
	Pid        int
	Ruid       int // real
	Euid       int // effective
	Suid       int // saved uid
	Fsuid      int // fs uid
	Gid        int
	Egid       int
	Sgid       int
	Fsgid      int
	CapInh     uint64
	CapPrm     uint64
	CapEff     uint64
	CapBnd     uint64
	CapAmb     uint64
	NoNewPrivs bool
	Seccomp    int
}

// ParseStatusFile parses a /proc/[pid]/status file and returns a Status struct.
func ParseStatusFile(filePath string) (*Status, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open status file %s: %w", filePath, err)
	}
	defer f.Close()
	return Parse(f)
}

// Parse parses the content of a /proc/[pid]/status file and returns a Status struct.
func Parse(r io.Reader) (*Status, error) {
	status := &Status{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Name":
			status.Name = value
		case "Pid":
			v, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse Pid: %w", err)
			}
			status.Pid = v
		case "Uid":
			uids, err := parseInts(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse Uid: %w", err)
			}
			if len(uids) >= 4 {
				status.Ruid = uids[0]
				status.Euid = uids[1]
				status.Suid = uids[2]
				status.Fsuid = uids[3]
			}
		case "Gid":
			gids, err := parseInts(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse Gid: %w", err)
			}
			if len(gids) >= 4 {
				status.Gid = gids[0]
				status.Egid = gids[1]
				status.Sgid = gids[2]
				status.Fsgid = gids[3]
			}
		case "CapInh":
			v, err := parseHex(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse CapInh: %w", err)
			}
			status.CapInh = v
		case "CapPrm":
			v, err := parseHex(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse CapPrm: %w", err)
			}
			status.CapPrm = v
		case "CapEff":
			v, err := parseHex(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse CapEff: %w", err)
			}
			status.CapEff = v
		case "CapBnd":
			v, err := parseHex(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse CapBnd: %w", err)
			}
			status.CapBnd = v
		case "CapAmb":
			v, err := parseHex(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse CapAmb: %w", err)
			}
			status.CapAmb = v
		case "NoNewPrivs":
			nnp, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse NoNewPrivs: %w", err)
			}
			status.NoNewPrivs = nnp != 0
		case "Seccomp":
			v, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse Seccomp: %w", err)
			}
			status.Seccomp = v
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return status, nil
}

func parseInts(s string) ([]int, error) {
	fields := strings.Fields(s)
	ints := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("failed to parse int %q: %w", f, err)
		}
		ints[i] = v
	}
	return ints, nil
}

func parseHex(s string) (uint64, error) {
	return strconv.ParseUint(s, 16, 64)
}
