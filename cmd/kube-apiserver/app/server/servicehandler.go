package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	kubeservice "minik8s/pkg/service"
	"net/http"
)

func UpdateService(c *gin.Context, s *Server) {
	prefix := "[serviceHandler][UpdateService]:"
	val, _ := io.ReadAll(c.Request.Body)
	service := kubeservice.Service{}
	err := json.Unmarshal([]byte(val), &service)
	if err != nil {
		fmt.Println(err)
		return
	}
	key := c.Request.URL.Path + "/" + service.ServiceMeta.Name
	body, _ := json.Marshal(service)
	fmt.Println(prefix + "key:" + key)
	fmt.Println(prefix + "value:" + string(body))
	err = s.Etcdstore.Put(key, string(body))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "etcd put service failed.",
		})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "update service success.",
	})
}

func DeleteService(c *gin.Context, s *Server) {
	prefix := "[serviceHandler][DeleteService]:"
	key := c.Request.URL.Path
	fmt.Println(prefix + "key:" + key)
	err := s.Etcdstore.Del(key)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delete service failed",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "delete service success",
	})
}
