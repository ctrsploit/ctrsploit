package containerd

import (
	"context"
	"fmt"
	"os"

	"github.com/containerd/containerd"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const pathContainerdSock = "/run/containerd/containerd.sock"

func GetVersionBySock() (v containerd.Version, err error) {
	if _, err := os.Stat(pathContainerdSock); err != nil {
		return v, fmt.Errorf("failed to stat %s: %w", pathContainerdSock, err)
	}
	client, err := containerd.New(pathContainerdSock)
	if err != nil {
		return v, fmt.Errorf("failed to create containerd client: %w", err)
	}
	defer client.Close()

	v, err = client.Version(context.Background())
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}
