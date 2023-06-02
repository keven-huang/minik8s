package main

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	"time"
)

func main() {
	node := &core.Node{}
	node.Name = "node2"
	err := tool.AddNode(node)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(1 * time.Second)
	Informer := informer.NewInformer("/api/v1/nodes")
	res := Informer.List()
	for _, v := range res {
		node := &core.Node{}
		err := json.Unmarshal([]byte(v.Value), &node)
		if err != nil {
			continue
		}
		fmt.Println(node.Name)
	}
}
