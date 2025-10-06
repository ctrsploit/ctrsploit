package rootfs

import (
	"fmt"
	"strings"

	"github.com/ctrsploit/ctrsploit/pkg/mountinfo"
	"github.com/ctrsploit/ctrsploit/pkg/runtime"
)

type ContainerdOverlay struct {
}

func (d ContainerdOverlay) Is() (bool, error) {
	// 1. not docker
	// docker+containerd: rootfs will be created by docker
	if is, _ := runtime.Docker().Is(); is {
		return false, nil
	}
	// 2. is containerd
	if is, _ := runtime.Containerd().Is(); !is {
		return false, nil
	}
	// 3. is overlay
	info, err := mountinfo.RootMount()
	if err != nil {
		return false, fmt.Errorf("error getting root's mount info: %w", err)
	}
	if !mountinfo.IsOverlay(info) {
		return false, nil
	}
	return true, nil
}

// RootPath
// k8s: /run/containerd/io.containerd.runtime.v2.task/k8s.io/6d14cd138fd87c7a4df37a00c38786d7c48db4114e46bad51ecb25f3174c1afa/rootfs/
// nerdctl: /run/containerd/io.containerd.runtime.v2.task/default/d7d37d49b063f3db3658cc3b674556339108cdabea6cc7801773e447ee01b40a/rootfs/
func (d ContainerdOverlay) RootPath() (string, error) {
	// TODO: try to get namespace, container-id
	// get container-id failed, return upperDir
	return d.upperDir()
}

func (d ContainerdOverlay) upperDir() (string, error) {
	info, err := mountinfo.RootMount()
	if err != nil {
		return "", fmt.Errorf("error getting root's mount info: %w", err)
	}
	const upperDirPrefix = "upperdir="
	options := strings.Split(info.VFSOptions, ",")
	for _, opt := range options {
		if strings.HasPrefix(opt, upperDirPrefix) {
			path := strings.TrimPrefix(opt, upperDirPrefix)
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 'upperdir' in overlay mount options: %s", info.VFSOptions)
}
