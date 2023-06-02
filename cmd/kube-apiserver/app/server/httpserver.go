package server

import (
	"log"
	"minik8s/pkg/kube-apiserver/etcd"

	"github.com/gin-gonic/gin"
)

type Server struct {
	// Api routes
	Router *gin.Engine
	// Etcd
	Etcdstore *etcd.EtcdStore
}

func NewServer() *Server {
	// initialize restful server
	router := gin.Default()
	// initialize etcd
	etcdstore, err := etcd.InitEtcdStore()
	if err != nil {
		log.Fatal(err)
	}
	s := &Server{Router: router, Etcdstore: etcdstore}
	return s
}

func (s *Server) RegisterHandler(handlers []Handler) {
	for _, handler := range handlers {
		h := handler
		switch h.Type {
		case "GET":
			s.Router.GET(h.Path, func(ctx *gin.Context) {
				h.Eventhandler(ctx, s)
			})
		case "PUT":
			s.Router.PUT(h.Path, func(ctx *gin.Context) {
				h.Eventhandler(ctx, s)
			})
		case "POST":
			s.Router.POST(h.Path, func(ctx *gin.Context) {
				h.Eventhandler(ctx, s)
			})
		case "DELETE":
			s.Router.DELETE(h.Path, func(ctx *gin.Context) {
				h.Eventhandler(ctx, s)
			})
		}
	}
}

func (s *Server) Run() {
	s.RegisterHandler(HandlerTable)
	err := s.Router.Run(":8080")
	if err != nil {
		return
	}
}
