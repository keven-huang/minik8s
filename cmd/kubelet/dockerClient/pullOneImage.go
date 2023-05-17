package main

import (
	"minik8s/pkg/kubelet/config"
	"minik8s/pkg/kubelet/dockerClient"
)

func main() {

	//image := "alpine"
	image := config.PAUSE_IMAGE_NAME

	err := dockerClient.PullOneImage(image)
	if err != nil {
		panic(err.Error())
	}
}
