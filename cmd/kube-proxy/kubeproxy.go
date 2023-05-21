package main

import kubeproxy "minik8s/pkg/kube-proxy"

func main() {
	proxy := kubeproxy.NewKubeProxy()
	proxy.Register()
	proxy.Run()
	select {}
}
