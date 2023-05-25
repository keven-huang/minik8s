package main

import (
	"fmt"
	jobcontroller "minik8s/pkg/kube-controller/job-controller"
	"minik8s/pkg/kube-controller/replicaset-controller"
)

func RunReplicaSetController() {
	ReplicaSetController, err := replicaset_controller.NewReplicaSetController()
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
