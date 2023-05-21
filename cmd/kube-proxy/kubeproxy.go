package main

import kubeproxy "minik8s/pkg/kube-proxy"

func main() {
	proxy := kubeproxy.NewKubeProxy()
	proxy.Run()
	select {}
}
