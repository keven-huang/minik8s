package main

import (
	"fmt"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
)

func watch_resource() {
	// initialize informer
	informer := informer.NewInformer("/api/v1/pods")
	informer.AddEventHandler(tool.Modified, func(event tool.Event) {
		// handle event
		fmt.Println("in handler")
		fmt.Println(event.Key)
	})
	informer.Run()
}

func main() {
	watch_resource()
}
