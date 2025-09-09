package selinux

import (
	_ "unsafe"

	_ "github.com/opencontainers/selinux/go-selinux"
)

//go:linkname GetSelinuxMountPoint github.com/opencontainers/selinux/go-selinux.getSelinuxMountPoint
func GetSelinuxMountPoint() string

//go:linkname FindSELinuxfs github.com/opencontainers/selinux/go-selinux.findSELinuxfs
func FindSELinuxfs() string

//go:linkname VerifySELinuxfsMount github.com/opencontainers/selinux/go-selinux.verifySELinuxfsMount
func VerifySELinuxfsMount(string) bool
