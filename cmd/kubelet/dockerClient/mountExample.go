package main

import (
	"github.com/docker/docker/api/types/volume"
	"minik8s/pkg/api/core"
	"minik8s/pkg/kubelet/dockerClient"
)

func VolumeExp() (volume.Volume, error) {
	return dockerClient.CreateVolume("E:/tmp/", "volume01")
}

func main() {
	resp, err := VolumeExp()
	if err != nil {
		panic(err.Error())
	}
	mt := core.VolumeMount{
		Name:      resp.Name,
		MountPath: "/tmp/mnt1/",
	}
	c1 := core.Container{
		Image:        "chasingdreams/minor_ubuntu:v1",
		Name:         "ubuntu_mnt_01",
		Command:      []string{"/bin/sh"},
		Tty:          true,
		VolumeMounts: []core.VolumeMount{mt},
	}
	resp2, err := dockerClient.CreateContainer(c1)
	if err != nil {
		panic(err.Error())
	} else {
		print(resp2.ID)
	}
}
