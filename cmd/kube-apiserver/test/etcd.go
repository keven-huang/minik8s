package main

import (
	"fmt"
	"log"
	"minik8s/pkg/kube-apiserver/etcd"
)

func main() {
	// initialize etcd
	fmt.Printf("===etcd test===\n")
	etcdstore, err := etcd.InitEtcdStore()
	if err != nil {
		log.Fatal(err)
	}
	etcdstore.PutKey("hello", "3")
	res, err := etcdstore.Get("hello")
	fmt.Printf("get  = %v", res[0])
}
