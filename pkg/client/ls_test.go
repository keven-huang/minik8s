package client

import (
	"fmt"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
)

func watch_resource() {
	// initialize informer
	informer := informer.NewInformer("/api/v1/pods/a")
	informer.AddEventHandler(func(event tool.Event) {
		// handle event
		fmt.Println(event)
	})
	informer.Run()
}
