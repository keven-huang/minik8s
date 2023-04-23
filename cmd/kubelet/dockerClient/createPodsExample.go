package main

import (
	"minik8s/pkg/api/core"
	"minik8s/pkg/kubelet"
	"minik8s/pkg/kubelet/dockerClient"
)

func main() {
	//alpine()
	//pause()
	//twoAlpine()
	ubuntu()
	//nginx()
}

func nginx() {
	c1 := core.Container{
		Image:   "nginx",
		Name:    "nginx",
		Command: []string{"/bin/sh"},
	}
	resp, err := dockerClient.CreateContainer(c1)
	if err != nil {
		panic("nginx")
	} else {
		print(resp.ID)
		print(resp.Warnings)
	}
}

func ubuntu() {
	c1 := core.Container{
		Image:      "ubuntu",
		Name:       "ubuntu_test01",
		EntryPoint: []string{"echo", "hello"},
	}
	resp, err := dockerClient.CreateContainer(c1)
	if err != nil {
		panic("ubuntu")
	} else {
		print(resp.ID)
		print(resp.Warnings)
	}
}

func twoAlpine() {
	c1 := core.Container{
		Image:   "alpine",
		Name:    "alpine_01",
		Command: []string{"/bin/sh"},
	}
	c2 := core.Container{
		Image:   "alpine",
		Name:    "alpine_02",
		Command: []string{"/bin/sh"},
	}
	cons := []core.Container{c1, c2}
	metas, _, err := dockerClient.CreatePod(cons)
	if err != nil {
		panic(err.Error())
	} else {
		print(metas)
	}
}

func pause() {
	var empty []kubelet.Port
	resp, err := dockerClient.CreatePauseContainer("test1_pause", empty)
	if err != nil {
		panic("fun")
	} else {
		print(resp.ID)
		print(resp.Warnings)
	}
}

func alpine() {
	name := "test01-alpine"
	image := "alpine"
	resp, err := dockerClient.CreateContainer(core.Container{
		Image: image,
		Name:  name,
	})
	if err != nil {
		panic("alpine")
	} else {
		print(resp.ID)
		print(resp.Warnings)
	}
}
