package server

import (
	"minik8s/cmd/kube-apiserver/app/apiconfig"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Type         string
	Path         string
	Eventhandler func(c *gin.Context, s *Server)
}

var HandlerTable = []Handler{
	// etcd
	{Type: "GET", Path: apiconfig.PATH, Eventhandler: Get},
	{Type: "PUT", Path: apiconfig.PATH, Eventhandler: Put},
	{Type: "POST", Path: apiconfig.PATH, Eventhandler: Post},
	{Type: "DELETE", Path: apiconfig.PATH, Eventhandler: Delete},
	// node
	{Type: "GET", Path: apiconfig.NODE_PATH, Eventhandler: GetNode},
	{Type: "PUT", Path: apiconfig.NODE_PATH, Eventhandler: AddNode},
	// pod
	{Type: "GET", Path: apiconfig.POD_PATH, Eventhandler: GetPod},
	{Type: "PUT", Path: apiconfig.POD_PATH, Eventhandler: AddPod},
	{Type: "DELETE", Path: apiconfig.POD_PATH, Eventhandler: DeletePod},
	{Type: "POST", Path: apiconfig.POD_PATH, Eventhandler: UpdatePod},
	{Type: "GET", Path: "/watch/*resource", Eventhandler: Watch},
	{Type: "POST", Path: apiconfig.JOB_PATH, Eventhandler: AddJob},
	{Type: "POST", Path: apiconfig.JOB_FILE_PATH, Eventhandler: AddJobFile},
	{Type: "GET", Path: apiconfig.JOB_FILE_PATH, Eventhandler: GetJobFile}}
