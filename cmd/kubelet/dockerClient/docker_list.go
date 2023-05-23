package main

import (
	"fmt"
	"minik8s/pkg/kubelet/dockerClient"
)

func main() {

	containers, err := dockerClient.GetAllContainers()
	if err != nil {
		panic(err.Error())
	}
	for _, con := range containers {
		fmt.Printf("%v %s [%s]\n", con.Names, con.ID, con.Status)
	}
}
