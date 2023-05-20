package main

import jobcontroller "minik8s/pkg/kube-controller/job-controller"

func main() {
	jobconroller := jobcontroller.NewJobController()
	jobconroller.Register()
	jobconroller.Run()
	select {}
}
