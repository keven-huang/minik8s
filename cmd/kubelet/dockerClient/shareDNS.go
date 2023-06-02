package main

import (
	"fmt"
	"minik8s/pkg/api/core"
	"minik8s/pkg/kubelet/dockerClient"
)

func main() {
	con := core.Container{
		Image:   "chasingdreams/minor_ubuntu:v3",
		Name:    "user",
		Tty:     true,
		Command: []string{"/bin/sh"},
	}
	container, err := dockerClient.CreateContainer(con)
	if err != nil {
		return
	}
	fmt.Println(container.ID)
}
