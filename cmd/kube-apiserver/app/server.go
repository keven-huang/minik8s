package app

import "minik8s/cmd/kube-apiserver/app/server"

func InitializeServer() *server.Server {
	return server.NewServer()
}
