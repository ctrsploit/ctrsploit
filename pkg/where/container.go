package where

import "github.com/ctrsploit/ctrsploit/pkg/namespace"

type Container struct {
}

// IsIn
// We believe that being in a child mount namespace is equivalent to being inside a container.
// This holds true even if you are in an unshare environment or a chroot environment.
func (c Container) IsIn() (in bool, err error) {
	arbitrator, err := namespace.NewInoArbitrator()
	if err != nil {
		return
	}
	level, err := namespace.GetNamespaceLevel(arbitrator, namespace.NameMnt)
	if err != nil {
		return
	}
	if level == namespace.LevelChild {
		in = true
	}
	return
}
