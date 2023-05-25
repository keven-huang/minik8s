package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	v1 "minik8s/pkg/apis/meta/v1"
	"minik8s/pkg/kube-apiserver/etcd"
	"minik8s/pkg/util/random"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddPod(c *gin.Context, s *Server) {
	fmt.Println("In AddPod.")
	val, _ := io.ReadAll(c.Request.Body)
	pod := core.Pod{}
	err := json.Unmarshal([]byte(val), &pod)
	if err != nil {
		log.Println(err)
		return
	}
	key := c.Request.URL.Path + "/" + pod.Name
	res, err := s.Etcdstore.Get(key)
	if len(res) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Pod Name Duplicate.",
		})
		return
	}

	// pod data
	pod.Status.Phase = "Pending"
	pod.ObjectMeta.CreationTimestamp = v1.Now()
	pod.ObjectMeta.Generation = 1
	pod.UID = random.GenerateUUID()

	body, _ := json.Marshal(pod)

	err = s.Etcdstore.Put(key, string(body))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "etcd put pod failed.",
		})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "add pod success.",
	})
}

// GetPod Body传入Pod.Name
func GetPod(c *gin.Context, s *Server) {
	fmt.Println("[api-server] [podHandler] [GetPod]")
	if c.Query("all") == "true" {
		// delete the keys
		res, err := s.Etcdstore.GetWithPrefix(apiconfig.POD_PATH)
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}

	PodName := c.Query("Name")
	key := c.Request.URL.Path + "/" + string(PodName)

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
			"message": "etcd get pod failed",
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func DeletePod(c *gin.Context, s *Server) {
	fmt.Println("In DeletePod")
	err := c.Request.ParseForm()
	if err != nil {
		return
	}
	if c.Query("all") == "true" {
		// delete the keys
		num, err := s.Etcdstore.DelAll(apiconfig.POD_PATH)
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":   "delete all pods successfully.",
			"deleteNum": num,
		})
		return
	}

	PodName := c.Query("PodName")
	fmt.Println("PodName:", PodName)
	key := c.Request.URL.Path + "/" + PodName
	err = s.Etcdstore.Del(key)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delete pod failed",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "delete pod success",
		"deletePodName": PodName,
	})
}

func UpdatePod(c *gin.Context, s *Server) {
	val, _ := io.ReadAll(c.Request.Body)
	pod := core.Pod{}
	err := json.Unmarshal([]byte(val), &pod)
	if err != nil {
		fmt.Println(err)
		return
	}
	key := c.Request.URL.Path + "/" + pod.Name

	body, _ := json.Marshal(pod)
	err = s.Etcdstore.Put(key, string(body))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "etcd put pod failed.",
		})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "update pod success.",
	})
}
