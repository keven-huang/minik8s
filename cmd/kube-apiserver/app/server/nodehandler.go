package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
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
		c.JSON(http.StatusOK, res)
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
	c.JSON(http.StatusOK, res)
}

func AddNode(c *gin.Context, s *Server) {
	val, _ := io.ReadAll(c.Request.Body)
	node := core.Node{}
	err := json.Unmarshal([]byte(val), &node)
	if err != nil {
		log.Println(err)
		return
	}
	key := c.Request.URL.Path + "/" + node.Name
	res, _ := s.Etcdstore.Get(key)
	if len(res) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Node Name Duplicate.",
		})
		return
	}
	// node data
	err = s.Etcdstore.Put(key, string(val))
	if err != nil {
		fmt.Print(err)
		return
	}
}
