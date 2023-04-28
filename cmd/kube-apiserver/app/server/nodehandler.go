package server

import (
	"log"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNode(c *gin.Context, s *Server) {
	if c.Query("all") == "true" {
		// delete the keys
		res, err := s.Etcdstore.GetAll(apiconfig.NODE_PATH)
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "get all node successfully.",
			"nodes":   res,
		})
		return
	}
	NodeName := c.Query("NodeName")
	key := c.Request.URL.Path + "/" + string(NodeName)
	res, err := s.Etcdstore.Get(key)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "etcd get node failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "get node successfully.",
		"Nodes":   res,
	})
}
