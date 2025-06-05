package containerd

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

func GetVersionBySock() (v containerd.Version, err error) {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	defer client.Close()

	v, err = client.Version(context.Background())
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	return
}
