package main

import (
	"fmt"
	"minik8s/pkg/api/core"
	metav1 "minik8s/pkg/apis/meta/v1"
	"minik8s/pkg/kubelet/dockerClient"

	"github.com/docker/docker/api/types/volume"
)

func VolumeExp() (volume.Volume, error) {
	// str is host-path
	str := "/tmp/host_path"
	return dockerClient.CreateVolume("volume02", &str)
}

// c1 and c2 share volume01
func example() {
	mt := core.VolumeMount{
		Name:      "volume03",
		MountPath: "/tmp/host_path",
	}
	con := core.Container{
		Image:        "chasingdreams/minor_ubuntu:v2",
		Name:         "test_mnt",
		EntryPoint:   []string{"/tmp/host_path/host_file"},
		Tty:          true,
		VolumeMounts: []core.VolumeMount{mt},
	}
	volume1 := core.Volume{Name: "volume03", HostPath: "/tmp/host_path"}
	pod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: con.Name,
		},
		Spec: core.PodSpec{
			Volumes: []core.Volume{
				volume1,
			},
			Containers: []core.Container{
				con,
			},
		},
	}
	_, _, err := dockerClient.CreatePod(*pod)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// func main() {
// 	example()
//resp, err := VolumeExp()
//if err != nil {
//	panic(err.Error())
//}
//
//mt := core.VolumeMount{
//	Name:      resp.Name,
//	MountPath: "/tmp/mnt1/",
//}
//mt2 := core.VolumeMount{
//	Name:      resp.Name,
//	MountPath: "/tmp/mnt2",
//}
//c2 := core.Container{
//	Image:        "chasingdreams/minor_ubuntu:v1",
//	Name:         "ubuntu_mnt_02",
//	Command:      []string{"/bin/sh"},
//	Tty:          true,
//	VolumeMounts: []core.VolumeMount{mt2},
//}
//c1 := core.Container{
//	Image:        "chasingdreams/minor_ubuntu:v1",
//	Name:         "ubuntu_mnt_01",
//	Command:      []string{"/bin/sh"},
//	Tty:          true,
//	VolumeMounts: []core.VolumeMount{mt},
//}
//resp2, err := dockerClient.CreateContainer(c1)
//if err != nil {
//	panic(err.Error())
//} else {
//	print(resp2.ID)
//}
//resp3, err2 := dockerClient.CreateContainer(c2)
//if err2 != nil {
//	panic(err2.Error())
//} else {
//	print(resp3.ID)
//}

// }
