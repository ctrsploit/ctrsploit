package sysctl

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	PathUserMaxUserNamespaces = "/proc/sys/user/max_user_namespaces"
	UserNamespaceNotSupported = -1
)

/*
MaxUserNamespaces returns -1 if file not exists, which means user namespace not supported
*/
func MaxUserNamespaces() (int, error) {
	content, err := os.ReadFile(PathUserMaxUserNamespaces)
	if err != nil {
		if os.IsNotExist(err) {
			return UserNamespaceNotSupported, nil
		}
		return 0, fmt.Errorf("failed to read user.max_user_namespaces: %w", err)
	}
	i, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		return 0, fmt.Errorf("failed to parse user.max_user_namespaces: %w", err)
	}
	return i, nil
}
