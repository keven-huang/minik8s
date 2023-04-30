package main

import (
	"minik8s/cmd/kube-apiserver/app"
)

func main() {
	s := app.InitializeServer()
	s.Run()
}
