package main

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	"minik8s/pkg/kubelet/dockerClient"
)

func CreatePod(event tool.Event) {
	// handle event
	fmt.Println("In AddPod EventHandler:")
	fmt.Println("event.Key: ", event.Key)
	fmt.Println("event.Val: ", event.Val)

	pod := &core.Pod{}
	err := json.Unmarshal([]byte(event.Val), pod)
	if err != nil {
		fmt.Println(err)
		return
	}

	metaData, netSetting, err := dockerClient.CreatePod(*pod)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, meta := range metaData {
		fmt.Println(meta.Name, meta.Id)
	}
	net, err := json.MarshalIndent(netSetting, "", "  ")
	fmt.Println("net: ")
	fmt.Println(string(net))
	fmt.Println("-----------")
}

func DeletePod(event tool.Event) {
	// handle event
	fmt.Println("In DeletePod EventHandler:")
	fmt.Println("event.Key: ", event.Key)

	pod := &core.Pod{}
	err := json.Unmarshal([]byte(event.Val), pod)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = dockerClient.DeletePod(*pod)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	Informer := informer.NewInformer("/api/v1/pods")

	// Create Pod.
	Informer.AddEventHandler(tool.Added, CreatePod)

	Informer.AddEventHandler(tool.Deleted, DeletePod)

	Informer.Run()
}
