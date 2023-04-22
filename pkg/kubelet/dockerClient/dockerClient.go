/*
	this client is based on docker-go-client

*/

package dockerClient

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
)

func GetNewClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.WithVersion("1.41"))
}

// GetAllContainers get all containers, including running / stopped
// docker ps
func GetAllContainers() ([]types.Container, error) {
	cli, err := GetNewClient()
	if err != nil {
		return nil, err
	}
	return cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
}

// GetRunningContainers get running containers
// docker ps -a
func GetRunningContainers() ([]types.Container, error) {
	cli, err := GetNewClient()
	if err != nil {
		return nil, err
	}
	return cli.ContainerList(context.Background(), types.ContainerListOptions{})
}

func GetAllImages() ([]types.ImageSummary, error) {
	cli, err := GetNewClient()
	if err != nil {
		return nil, err
	}
	return cli.ImageList(context.Background(), types.ImageListOptions{All: true})
}

// ImageInSet judge if image is in a set
// return true if in,
func ImageInSet(image string, imageSet []string) bool {
	fullName := image + ":latest" // default is latest
	for _, name := range imageSet {
		if name == image {
			return true
		}
		if fullName == name {
			return true
		}
	}
	print("Image doesn't exist: %v", image)
	return false
}

// PullImages  pull some images in the array
// ignore the images which already exists in local
func PullImages(targets []string) error {
	historyImages, err := GetAllImages()
	if err != nil {
		return err
	}
	var filteredTarget []string
	for _, image := range targets {
		var flag bool
		for _, historyImage := range historyImages {
			if ImageInSet(image, historyImage.RepoTags) {
				flag = true
				break
			}
		}
		if flag == false {
			filteredTarget = append(filteredTarget, image)
		}
	}
	if filteredTarget != nil {
		for _, target := range filteredTarget {
			err := PullOneImage(target)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// PullOneImage 拉取一个镜像
// ImagePull自带重试功能
// 如果没有找到，会报panic
func PullOneImage(target string) error {
	print("start to pull image: ", target)
	cli, err := GetNewClient()
	if err != nil {
		return err
	}
	out, err := cli.ImagePull(context.Background(), target, types.ImagePullOptions{})
	if err != nil {
		panic(err.Error())
		return err
	}
	// 延迟返回，在这个pullOneImage函数完全结束之后延迟返回
	defer out.Close()
	io.Copy(io.Discard, out)
	return nil
}

// 根据containerId start一个已有的container
// 有些容器可能无法使用
func StartOneContainer(id string) error {
	cli, err := GetNewClient()
	if err != nil {
		return err
	}
	err = cli.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

// 如果容器没有启动，等价于start
func RestartContainer(id string) error {
	cli, err := GetNewClient()
	if err != nil {
		return err
	}
	err = cli.ContainerRestart(context.Background(), id, container.StopOptions{})
	if err != nil {
		return err
	}
	return nil
}

// KillContainer 向id-container发送一个signal
// 比如SIGKILL, SIGINT
func KillContainer(id string, signal string) error {
	cli, err := GetNewClient()
	if err != nil {
		return err
	}
	err = cli.ContainerKill(context.Background(), id, signal)
	if err != nil {
		return err
	}
	return nil
}

func GetContainerInfo(id string) (types.ContainerJSON, error) {
	cli, err := GetNewClient()
	if err != nil {
		return types.ContainerJSON{}, err
	}
	info, err := cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return types.ContainerJSON{}, err
	}
	return info, nil
}
