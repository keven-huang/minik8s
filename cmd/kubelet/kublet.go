package main

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	"minik8s/pkg/kubelet/dockerClient"
)

func main() {
	Informer := informer.NewInformer("/api/v1/pods")
	Informer.AddEventHandler(tool.Added, func(event tool.Event) {
		// handle event
		fmt.Println("in handler")
		fmt.Println(event.Key)

		pod := &core.Pod{}
		err := json.Unmarshal([]byte(event.Val), pod)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, container := range pod.Spec.Containers {
			resp, err := dockerClient.CreateContainer(container)
			if err != nil {
				panic("Pod: " + pod.Name + "container: " + container.Name)
			} else {
				fmt.Println("Container ID: ", resp.ID)
				fmt.Println("Warn: ", resp.Warnings)
			}
		}

	})
	Informer.Run()
}
