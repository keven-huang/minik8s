package main

import (
	"minik8s/pkg/kubelet"
	"minik8s/pkg/kubelet/dockerClient"
)

func main() {

	//image := "alpine"
	image := kubelet.PAUSE_IMAGE_NAME

	err := dockerClient.PullOneImage(image)
	if err != nil {
		panic(err.Error())
	}
}
