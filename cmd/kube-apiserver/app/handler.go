package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"minik8s/pkg/api/core"
	v1 "minik8s/pkg/apis/meta/v1"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Put(c *gin.Context) {
	fmt.Printf("put\n")
	key := c.Request.URL.Path
	value, _ := io.ReadAll(c.Request.Body)
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
		log.Println(err)
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
	value, _ := io.ReadAll(c.Request.Body)
	err := s.etcdstore.Put(key, string(value))
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "post success",
	})
}

func (s *Server) AddPod(c *gin.Context) {
	val, _ := io.ReadAll(c.Request.Body)
	pod := core.Pod{}
	err := json.Unmarshal([]byte(val), &pod)
	if err != nil {
		log.Println(err)
		return
	}
	key := c.Request.URL.Path + "/" + pod.Name
	res, err := s.etcdstore.Get(key)
	if len(res) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Pod Name Duplicate.",
		})
	}

	// pod data
	pod.Status.Phase = "Pending"
	pod.ObjectMeta.CreationTimestamp = v1.Now()
	pod.ObjectMeta.Generation = 1

	body, _ := json.Marshal(pod)

	err = s.etcdstore.Put(key, string(body))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "etcd put pod failed.",
		})
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "add pod success.",
	})
}

// GetPod Body传入Pod.Name
func (s *Server) GetPod(c *gin.Context) {
	val, _ := io.ReadAll(c.Request.Body)
	key := c.Request.URL.Path + "/" + string(val)
	res, err := s.etcdstore.Get(key)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "etcd get pod failed",
		})
	}
	pod := core.Pod{}
	err = json.Unmarshal([]byte(res[0]), &pod)
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": pod,
	})
}

func (s *Server) DeletePod(c *gin.Context) {
	val, _ := io.ReadAll(c.Request.Body)
	key := c.Request.URL.Path + "/" + string(val)
	err := s.etcdstore.Del(key)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delete pod failed",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "delete pod success",
	})
}
