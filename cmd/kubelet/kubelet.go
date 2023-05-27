package main

import (
	"flag"
	"fmt"
	"minik8s/pkg/kubelet"
)

func main() {
	NodeName := flag.String("nodename", "node1", "node name")
	NodeIP := flag.String("nodeip", "127.0.0.1", "node ip")
	MasterIP := flag.String("masterip", "", "master ip")
	flag.Parse()
	k, err := kubelet.NewKubelet(*NodeName, *NodeIP, *MasterIP)
	if err != nil {
		fmt.Println(err)
		return
	}
	k.Register()
	k.Run()
}
