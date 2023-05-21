package main

import (
	"github.com/docker/docker/api/types/volume"
	"minik8s/pkg/api/core"
	"minik8s/pkg/kubelet/dockerClient"
)

func VolumeExp() (volume.Volume, error) {
	// str is host-path
	str := "/tmp/host_path"
	return dockerClient.CreateVolume("volume02", &str)
}

// c1 and c2 share volume01

func main() {
	resp, err := VolumeExp()
	if err != nil {
		panic(err.Error())
	}

	mt := core.VolumeMount{
		Name:      resp.Name,
		MountPath: "/tmp/mnt1/",
	}
	mt2 := core.VolumeMount{
		Name:      resp.Name,
		MountPath: "/tmp/mnt2",
	}
	c2 := core.Container{
		Image:        "chasingdreams/minor_ubuntu:v1",
		Name:         "ubuntu_mnt_02",
		Command:      []string{"/bin/sh"},
		Tty:          true,
		VolumeMounts: []core.VolumeMount{mt2},
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
	resp3, err2 := dockerClient.CreateContainer(c2)
	if err2 != nil {
		panic(err2.Error())
	} else {
		print(resp3.ID)
	}

}
