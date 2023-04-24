package main

import (
	"minik8s/pkg/api/core"
	"minik8s/pkg/kubelet/dockerClient"
)

func main() {
	//alpine()
	//pause()
	twoUbuntu()
	//minorUbuntu()
	//ubuntu()
	//nginx()
}

func minorUbuntu() {
	c1 := core.Container{
		Image:   "minor_ubuntu:v1",
		Name:    "ubuntu_01",
		Command: []string{"/bin/sh"},
		Tty:     true,
	}
	c2 := core.Container{
		Image:   "minor_ubuntu:v1",
		Name:    "ubuntu_02",
		Command: []string{"/bin/sh"},
		Tty:     true,
	}
	cons := []core.Container{c1, c2}
	metas, _, err := dockerClient.CreatePod(cons)
	if err != nil {
		panic(err.Error())
	} else {
		print(metas)
	}
}

func nginx() {
	c1 := core.Container{
		Image:      "nginx",
		Name:       "nginx",
		EntryPoint: []string{"/bin/sh"},
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
		EntryPoint: []string{"sh"},
	}
	resp, err := dockerClient.CreateContainer(c1)
	if err != nil {
		panic("ubuntu")
	} else {
		print(resp.ID)
		print(resp.Warnings)
	}
}

func twoUbuntu() {
	c1 := core.Container{
		Image:   "chasingdreams/minor_ubuntu:v1",
		Name:    "ubuntu_01",
		Command: []string{"/bin/sh"},
		Tty:     true,
	}
	c2 := core.Container{
		Image:   "chasingdreams/minor_ubuntu:v1",
		Name:    "ubuntu_02",
		Command: []string{"/bin/sh"},
		Tty:     true,
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
	var empty []core.Port
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
