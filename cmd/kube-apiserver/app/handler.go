package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"minik8s/pkg/api/core"
	v1 "minik8s/pkg/apis/meta/v1"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Put(c *gin.Context) {
	fmt.Printf("put\n")
	key := c.Request.URL.Path
	value, _ := ioutil.ReadAll(c.Request.Body)
	err := s.etcdstore.Put(key, string(value))
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "put success",
	})
}

func (s *Server) Get(c *gin.Context) {
	key := c.Request.URL.Path
	res, err := s.etcdstore.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
}

func (s *Server) Delete(c *gin.Context) {
	key := c.Request.URL.Path
	err := s.etcdstore.Del(key)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "delete success",
	})
}

func (s *Server) POST(c *gin.Context) {
	key := c.Request.URL.Path
	value, _ := ioutil.ReadAll(c.Request.Body)
	err := s.etcdstore.Put(key, string(value))
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "post success",
	})
}

func (s *Server) AddPod(c *gin.Context) {
	key := c.Request.URL.Path
	val, _ := ioutil.ReadAll(c.Request.Body)
	pod := core.Pod{}
	json.Unmarshal([]byte(val), &pod)
	body, _ := json.Marshal(pod)

	// pod data
	pod.Status.Phase = "Pending"
	pod.ObjectMeta.CreationTimestamp = v1.Now()
	pod.ObjectMeta.Generation = 0

	err := s.etcdstore.Put(key, string(body))
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "add pod failed",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "add pod success",
	})
}

func (s *Server) GetPod(c *gin.Context) {
	key := c.Request.URL.Path
	res, err := s.etcdstore.Get(key)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get pod failed",
		})
	}
	pod := core.Pod{}
	json.Unmarshal([]byte(res[0]), &pod)
	c.JSON(http.StatusOK, gin.H{
		"message": pod,
	})
}

func (s *Server) DeletePod(c *gin.Context) {
	key := c.Request.URL.Path
	err := s.etcdstore.Del(key)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delete pod failed",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "delete pod success",
	})
}
