package main

import (
	"fmt"
	"minik8s/pkg/controller/replicaset"
	jobcontroller "minik8s/pkg/kube-controller/job-controller"
)

func RunReplicaSetController() {
	ReplicaSetController, err := replicaset.NewReplicaSetController()
	if err != nil {
		fmt.Println(err)
		return
	}
	ReplicaSetController.Register()
	ReplicaSetController.Run()
}

func RunJobController() {
	jobconroller := jobcontroller.NewJobController()
	jobconroller.Register()
	jobconroller.Run()
}

func main() {
	go RunReplicaSetController()
	go RunJobController()
	select {}
}
