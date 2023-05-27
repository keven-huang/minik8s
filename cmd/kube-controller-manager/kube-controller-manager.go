package main

import (
	"fmt"
	hpacontroller "minik8s/pkg/kube-controller/hpa-controller"
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

func RunHpaController() {
	hpaController := hpacontroller.NewHPAController()
	hpaController.Register()
	hpaController.Run()
}

func main() {
	go RunReplicaSetController()
	go RunJobController()
	go RunHpaController()
	select {}
}
