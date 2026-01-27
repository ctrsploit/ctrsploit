package sysctl

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	PatternRouteLocalNet = "/proc/sys/net/ipv4/conf/*/route_localnet"
)

/*
RouteLocalNetEnabled
route_localnet - BOOLEAN

	Do not consider loopback addresses as martian source or destination
	while routing. This enables the use of 127/8 for local routing purposes.
	default FALSE

https://github.com/kubernetes/kubernetes/issues/90259
*/
func RouteLocalNetEnabled() (bool, error) {
	filePaths, err := filepath.Glob(PatternRouteLocalNet)
	if err != nil {
		return false, fmt.Errorf("failed to glob route_localnet: %w", err)
	}
	if len(filePaths) == 0 {
		return false, fmt.Errorf("no route_localnet found")
	}
	for _, path := range filePaths {
		content, err := os.ReadFile(path)
		if err != nil {
			return false, fmt.Errorf("failed to read %s: %w", path, err)
		}
		if strings.TrimSpace(string(content)) == "1" {
			return true, nil
		}
	}
	return false, nil
}

// RpFilterMode reads all rp_filter settings from /proc/sys/net/ipv4/conf/*/rp_filter
//
// rp_filter – INTEGER
// 0 - No source validation.
// 1 - Strict mode: each incoming packet is tested against the FIB entry that would be used
//
//	to send packets to the source address. If the interface does not match, the packet is dropped.
//
// 2 - Loose mode: packet is accepted if the source address is reachable via any interface.
// Default is 0 (off).
//
// Returns a map[string]int like {"all": 1, "eth0": 2, "lo": 0}
func RpFilterMode() (map[string]int, error) {
	result := make(map[string]int)

	baseDir := "/proc/sys/net/ipv4/conf"
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// skip unreadable dirs
			return nil
		}

		// Only look for rp_filter files directly under interface directories
		if info.Name() == "rp_filter" {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", path, err)
			}

			s := strings.TrimSpace(string(data))
			mode, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("invalid rp_filter value in %s: %q", path, s)
			}

			iface := filepath.Base(filepath.Dir(path))
			result[iface] = mode
		}
		return nil
	})

	return result, err
}
