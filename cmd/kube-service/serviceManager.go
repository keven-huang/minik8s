package kube_service

import kube_service "minik8s/pkg/runtime/kube-service"

func main() {
	sm := kube_service.NewServiceManager()
	sm.Run()
}
