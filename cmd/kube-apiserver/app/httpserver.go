package app

import (
	"log"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/kube-apiserver/etcd"

	"github.com/gin-gonic/gin"
)

// k8s api-server分为三层: 1. restful api 2. 鉴权 3. Registry

// TODO : registry scheme
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
		router.PUT(apiconfig.PATH, s.Put)
		router.POST(apiconfig.PATH, s.POST)
		router.DELETE(apiconfig.PATH, s.Delete)
	}
	// Pod Handler
	{
		router.GET(apiconfig.POD_PATH, s.GetPod)
		router.PUT(apiconfig.POD_PATH, s.AddPod)
		router.DELETE(apiconfig.POD_PATH, s.DeletePod)
	}

	return s
}

func (s *Server) Run() {
	s.router.Run(":8080")
}
