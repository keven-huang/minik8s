package app

import (
	"log"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/kube-apiserver/etcd"

	"github.com/gin-gonic/gin"
)

// k8s api-server分为三层: 1. restful api 2. 鉴权 3. Registry

type Server struct {
	// Api routes
	router *gin.Engine
	// Etcd
	etcdstore *etcd.EtcdStore
}

func NewServer() *Server {
	// initialize restful server
	router := gin.Default()
	// initialize etcd
	etcdstore, err := etcd.InitEtcdStore()
	if err != nil {
		log.Fatal(err)
	}
	s := &Server{router: router, etcdstore: etcdstore}

	// api 配置
	{
		router.GET(apiconfig.PATH, s.Get)
		router.POST(apiconfig.PATH, s.Put)
		router.DELETE(apiconfig.PATH, s.Delete)
		router.PATCH(apiconfig.PATH, s.Update)
	}

	return s
}

func (s *Server) Run() {
	s.router.Run(":8080")
}
