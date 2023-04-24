/*
	this client is based on docker-go-client

*/

package dockerClient

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"minik8s/pkg/api/core"
	"minik8s/pkg/kubelet"
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
	m := make(map[string]int)
	for _, image := range targets {
		var flag bool
		for _, historyImage := range historyImages {
			if ImageInSet(image, historyImage.RepoTags) {
				flag = true
				break
			}
		}
		_, ok := m[image]
		if ok == true {
			continue
		}
		if flag == false {
			filteredTarget = append(filteredTarget, image)
			fmt.Printf("Image doesn't exist: %v\n", image)
			m[image] = 1
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

func StopContainer(id string) error {
	cli, err := GetNewClient()
	if err != nil {
		return err
	}
	err = cli.ContainerStop(context.Background(), id, container.StopOptions{})
	return err
}

func DeleteContainer(id string) error {
	cli, err := GetNewClient()
	if err != nil {
		return err
	}
	_, err = cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return nil
	}
	err = cli.ContainerStop(context.Background(), id, container.StopOptions{})
	if err != nil {
		return err
	}
	err = cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

func DeleteContainers(ids []string) error {
	for _, name := range ids {
		err := DeleteContainer(name)
		if err != nil {
			return err
		}
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

/*
获取容器的网络设置
如果用在pause容器中，可以获取pod的网络设置
*/
func GetNetworkSettings(id string) (*types.NetworkSettings, error) {
	cli, err := GetNewClient()
	if err != nil {
		return nil, err
	}
	res, err := cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return res.NetworkSettings, nil
}

func CreateContainer(con core.Container) (container.CreateResponse, error) {
	cli, err := GetNewClient()
	if err != nil {
		return container.CreateResponse{}, err
	}
	err = PullImages([]string{con.Image})
	if err != nil {
		return container.CreateResponse{}, err
	}
	resp, err := cli.ContainerCreate(context.Background(),
		&container.Config{Image: con.Image, Entrypoint: con.EntryPoint, Tty: con.Tty},
		&container.HostConfig{}, nil, nil, con.Name)
	return resp, err
}

func CreatePauseContainer(name string, ports []core.Port) (container.CreateResponse, error) {
	cli, err := GetNewClient()
	if err != nil {
		return container.CreateResponse{}, err
	}
	// this for outside-call
	err = PullImages([]string{kubelet.PAUSE_IMAGE_NAME})
	if err != nil {
		return container.CreateResponse{}, err
	}
	var portSet nat.PortSet
	portSet = make(nat.PortSet, len(ports))
	for _, port := range ports {
		if port.Protocol == "tcp" {
			res, err := nat.NewPort("tcp", port.PortNumber)
			if err != nil {
				return container.CreateResponse{}, err
			}
			portSet[res] = struct{}{}
			continue
		}
		if port.Protocol == "udp" {
			res, err := nat.NewPort("udp", port.PortNumber)
			if err != nil {
				return container.CreateResponse{}, err
			}
			portSet[res] = struct{}{}
			continue
		}
		panic("unsupported network protocol")
	}
	response, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image:        kubelet.PAUSE_IMAGE_NAME,
		ExposedPorts: portSet,
	}, &container.HostConfig{
		IpcMode: container.IPCModeShareable,
		DNS:     []string{kubelet.DnsAddress},
	}, nil, nil, name)
	if err != nil {
		panic(err.Error())
	}
	return response, err
}

func CreatePod(containers []core.Container) ([]core.ContainerMeta, *types.NetworkSettings, error) {
	cli, err := GetNewClient()
	if err != nil {
		return nil, nil, err
	}
	var totalPorts []core.Port
	var res []core.ContainerMeta
	images := []string{kubelet.PAUSE_IMAGE_NAME}
	names := []string{kubelet.PAUSE_NAME}
	curPauseName := kubelet.PAUSE_NAME
	for _, v := range containers {
		curPauseName += "_" + v.Name
		names = append(names, v.Name)
		images = append(images, v.Image)
		for _, port := range v.Ports {
			totalPorts = append(totalPorts, port)
		}
	}
	// 首先删除之前同名的container
	err = DeleteContainers(names)
	if err != nil {
		return nil, nil, err
	}
	// 拉取镜像列表
	err = PullImages(images)
	if err != nil {
		return nil, nil, err
	}
	// 创建pause
	pauseContainer, err := CreatePauseContainer(curPauseName, totalPorts)
	if err != nil {
		return nil, nil, err
	}
	curPauseID := pauseContainer.ID
	res = append(res, core.ContainerMeta{Name: curPauseName, Id: curPauseID})
	for _, v := range containers {
		resp, err := cli.ContainerCreate(context.Background(), &container.Config{
			Image:      v.Image,
			Entrypoint: v.EntryPoint,
			Cmd:        v.Command,
			Tty:        v.Tty,
		}, &container.HostConfig{
			NetworkMode: container.NetworkMode("container:" + curPauseID),
			IpcMode:     container.IpcMode("container:" + curPauseID),
			PidMode:     container.PidMode("container:" + curPauseID),
		}, nil, nil, v.Name)
		if err != nil {
			return nil, nil, err
		}
		res = append(res, core.ContainerMeta{Name: v.Name, Id: resp.ID})
	}
	// 启动
	for _, v := range res {
		err := cli.ContainerStart(context.Background(), v.Id, types.ContainerStartOptions{})
		if err != nil {
			return nil, nil, err
		}
	}
	netSetting, err := GetNetworkSettings(curPauseID)
	if err != nil {
		return nil, nil, err
	}
	return res, netSetting, nil
}
