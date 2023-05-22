package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/kube-apiserver/etcd"
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
	if c.Query("all") == "true" { // delete all services
		num, err := s.Etcdstore.DelAll(apiconfig.SERVICE_PATH)
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":   "delete all services successfully.",
			"deleteNum": num,
		})
		return
	}
	ServiceName := c.Query("ServiceName")
	//fmt.Println("ServiceName:", ServiceName)
	key := c.Request.URL.Path + "/" + ServiceName
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
		"message":           "delete service success",
		"deleteServiceName": ServiceName,
	})
}

// GetService Body传入Service.Name
func GetService(c *gin.Context, s *Server) {
	fmt.Println("[api-server] [ServiceHandler] [GetService]")
	if c.Query("all") == "true" {
		res, err := s.Etcdstore.GetWithPrefix(apiconfig.SERVICE_PATH)
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}

	ServiceName := c.Query("Name")
	key := c.Request.URL.Path + "/" + string(ServiceName)

	var res []etcd.ListRes
	var err error

	if c.Query("prefix") == "true" {
		res, err = s.Etcdstore.GetWithPrefix(key)
		fmt.Println(res)
	} else {
		res, err = s.Etcdstore.GetExact(key)
	}

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "etcd get Service failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get Service successfully.",
		"Results": res,
	})
}
