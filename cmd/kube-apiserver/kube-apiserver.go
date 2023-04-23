package main

import (
	"minik8s/cmd/kube-apiserver/app"
)

func main() {
	server := app.NewServer()
	server.Run()
}
