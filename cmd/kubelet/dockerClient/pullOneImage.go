package main

import "minik8s/pkg/kubelet/dockerClient"

func main() {

	image := "alpine"

	err := dockerClient.PullOneImage(image)
	if err != nil {
		panic(err.Error())
	}
}
