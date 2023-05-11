package main

import "minik8s/pkg/kubelet"

func main() {
	k := kubelet.NewKublet()
	k.Register()
	k.Run()
}
