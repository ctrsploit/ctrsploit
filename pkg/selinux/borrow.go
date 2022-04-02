package selinux

import (
	_ "github.com/opencontainers/selinux/go-selinux"
	_ "unsafe"
)

//go:linkname GetSelinuxMountPoint github.com/opencontainers/selinux/go-selinux.getSelinuxMountPoint
func GetSelinuxMountPoint() string
