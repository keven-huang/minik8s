package main

import (
	"fmt"
	"minik8s/pkg/controller/replicaset"
)

func main() {
	ReplicaSetController, err := replicaset.NewReplicaSetController()
	if err != nil {
		fmt.Println(err)
		return
	}
	ReplicaSetController.Register()
	ReplicaSetController.Run()

	select {}
}
