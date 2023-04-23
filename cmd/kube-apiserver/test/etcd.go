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
	etcdstore.Put("hello", "3")
	res, err := etcdstore.Get("hello")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("res: %s", res[0])
}
