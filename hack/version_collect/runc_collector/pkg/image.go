package pkg

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"strings"
)

func inspect(image string) {
	cli, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// 使用 Docker 客户端调用 Inspect API
	imageInspect, _, err := cli.ImageInspectWithRaw(context.Background(), image)
	if err != nil {
		fmt.Printf("Error inspecting image: %s\n", err)
		return
	}
	lower := imageInspect.GraphDriver.Data["LowerDir"]
	dirs := strings.Split(lower, ":")
	fmt.Printf("Image info: %+v\n", dirs)

}

func GetFileHostPath(image string, filepath string) (hostPath string, err error) {
	return
}
