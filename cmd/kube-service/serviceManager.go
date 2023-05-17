package main

import kubeservice "minik8s/pkg/runtime/kube-service"

func main() {
	sm := kubeservice.NewServiceManager()
	sm.Run()
}
