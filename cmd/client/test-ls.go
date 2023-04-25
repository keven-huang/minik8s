package main

import (
	"fmt"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
)

func watchResource() {
	// initialize informer
	Informer := informer.NewInformer("/api/v1/pods")
	Informer.AddEventHandler(tool.Modified, func(event tool.Event) {
		// handle event
		fmt.Println("in modified handler")

	})
	Informer.Run()
}

func main() {
	watchResource()
}
