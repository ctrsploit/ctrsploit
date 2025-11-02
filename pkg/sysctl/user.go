package sysctl

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	PathUserMaxUserNamespaces = "/proc/sys/user/max_user_namespaces"
)

/*
MaxUserNamespaces returns 0 if file not exists, which means user namespace not supported
*/
func MaxUserNamespaces() (uint64, error) {
	content, err := os.ReadFile(PathUserMaxUserNamespaces)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	i, err := strconv.ParseUint(strings.TrimSpace(string(content)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse user.max_user_namespaces: %w", err)
	}
	return i, nil
}
