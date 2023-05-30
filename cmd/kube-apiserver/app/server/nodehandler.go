package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/util/random"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNode(c *gin.Context, s *Server) {
	if c.Query("all") == "true" {
		// delete the keys
		res, err := s.Etcdstore.GetWithPrefix(apiconfig.NODE_PATH)
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

func DeleteNode(c *gin.Context, s *Server) {
	prefix := "[apiserver][nodehandler]"
	nodeName := c.Query("Name")
	key := c.Request.URL.Path + "/" + string(nodeName)
	fmt.Println(prefix + ": delete node key:" + key)
	err := s.Etcdstore.Del(key)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delete node failed",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "delete node success",
		"key":     key,
	})
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
	//res, _ := s.Etcdstore.Get(key)

	// 重名就当做覆盖
	//if len(res) > 0 {
	//	c.JSON(web.StatusBadRequest, gin.H{
	//		"message": "Node Name Duplicate.",
	//	})
	//	return
	//}

	node.UID = random.GenerateUUID()

	body, _ := json.Marshal(node)

	err = s.Etcdstore.Put(key, string(body))
	if err != nil {
		fmt.Print(err)
		return
	}
}
