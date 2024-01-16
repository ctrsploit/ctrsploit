package pkg

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func inspectGraphDrivers(image string) (path []string, err error) {
	cli, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"), client.WithAPIVersionNegotiation())
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}

	imageInspect, _, err := cli.ImageInspectWithRaw(context.Background(), image)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	lower := imageInspect.GraphDriver.Data["LowerDir"]
	path = strings.Split(lower, ":")
	return
}

func GetFileHostPath(image string, path string) (hostPath string, err error) {
	dirs, err := inspectGraphDrivers(image)
	if err != nil {
		return
	}
	for _, dir := range dirs {
		_, err = os.Lstat(filepath.Join(dir, path))
		if !os.IsNotExist(err) {
			hostPath = filepath.Join(dir, path)
			err = nil
			return
		}
	}
	return
}

func ListTags(src string) (tags []string, err error) {
	ctx := context.Background()
	repo, err := name.NewRepository(src)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}

	puller, err := remote.NewPuller()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}

	lister, err := puller.Lister(ctx, repo)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}

	for lister.HasNext() {
		t, e := lister.Next(ctx)
		if e != nil {
			err = e
			awesome_error.CheckErr(err)
			return
		}
		tags = append(tags, t.Tags...)
	}
	return
}

func Pull(image string) (err error) {
	cli, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"), client.WithAPIVersionNegotiation())
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	_, _, err = cli.ImageInspectWithRaw(context.Background(), image)
	if err != nil {
		if strings.Contains(err.Error(), "No such image") {
			err = nil
		} else {
			awesome_error.CheckErr(err)
			return
		}
	} else {
		return
	}

	out, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	fmt.Printf("%s pulling\n", image)
	_, err = io.Copy(io.Discard, out)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s pulled\n", image)
	return
}
