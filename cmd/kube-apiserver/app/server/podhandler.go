package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	v1 "minik8s/pkg/apis/meta/v1"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddPod(c *gin.Context, s *Server) {
	fmt.Println("in add pod")
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

	if c.Query("all") == "true" {
		// delete the keys
		res, err := s.Etcdstore.GetAll(apiconfig.POD_PATH)
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "get all pods successfully.",
			"Pods":    res,
		})
		return
	}

	PodName := c.Query("PodName")
	key := c.Request.URL.Path + "/" + string(PodName)
	res, err := s.Etcdstore.Get(key)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "etcd get pod failed",
		})
		return
	}
	//pod := core.Pod{}
	//err = json.Unmarshal([]byte(res[0]), &pod)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	c.JSON(http.StatusOK, gin.H{
		"message": "get pod successfully.",
		"Pods":    res,
	})
}

func DeletePod(c *gin.Context, s *Server) {
	fmt.Println("in delete")
	err := c.Request.ParseForm()
	if err != nil {
		return
	}
	if c.Request.Form.Get("all") == "true" {
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

	PodName := c.Request.PostForm.Get("PodName")
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
