package cgroup

import "github.com/ctrsploit/ctrsploit/pkg/cgroup/version"

type Cgroup interface {
	GetVersion() version.Version
	IsTop(mountpoint, subsystemName string) (top bool, err error)
	ListSubsystems(mountpoint string) (subsystems []string, err error)
}
