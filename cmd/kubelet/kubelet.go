package main

import (
	"fmt"
	"minik8s/pkg/kubelet"
)

func main() {
	k, err := kubelet.NewKubelet("node1")
	if err != nil {
		fmt.Println(err)
		return
	}
	k.Register()
	k.Run()
}
