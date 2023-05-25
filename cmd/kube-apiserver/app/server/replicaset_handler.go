package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	v1 "minik8s/pkg/apis/meta/v1"
	"minik8s/pkg/kube-apiserver/etcd"
	"minik8s/pkg/util/random"
	"net/http"
)

func AddReplicaSet(c *gin.Context, s *Server) {
	prefix := "[api-server] [AddReplicaSet] "
	val, _ := io.ReadAll(c.Request.Body)
	r := core.ReplicaSet{}
	err := json.Unmarshal([]byte(val), &r)
	if err != nil {
		log.Println("[ERROR] ", prefix, err)
		return
	}
	key := c.Request.URL.Path + "/" + r.Name
	res, _ := s.Etcdstore.Get(key)
	if len(res) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Replica Name Duplicate.",
		})
		return
	}

	r.UID = random.GenerateUUID()
	r.ObjectMeta.CreationTimestamp = v1.Now()

	body, _ := json.Marshal(r)

	err = s.Etcdstore.Put(key, string(body))
	if err != nil {
		log.Println("[ERROR] ", prefix, err)
		return
	}
}

func UpdateReplicaSet(c *gin.Context, s *Server) {
	prefix := "[api-server] [UpdateReplicaSet] "
	val, _ := io.ReadAll(c.Request.Body)
	r := core.ReplicaSet{}
	err := json.Unmarshal([]byte(val), &r)
	if err != nil {
		log.Println("[ERROR] ", prefix, err)
		return
	}
	key := c.Request.URL.Path + "/" + r.Name

	body, _ := json.Marshal(r)

	err = s.Etcdstore.Put(key, string(body))
	if err != nil {
		log.Println("[ERROR] ", prefix, err)
		return
	}
}

// GetReplicaSet Body传入Name
func GetReplicaSet(c *gin.Context, s *Server) {
	prefix := "[api-server] [ReplicaSetHandler] [GetReplicaSet]"
	fmt.Println(prefix)
	if c.Query("all") == "true" {
		// delete the keys
		res, err := s.Etcdstore.GetWithPrefix(apiconfig.REPLICASET_PATH)
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}

	Name := c.Query("Name")
	key := c.Request.URL.Path + "/" + string(Name)

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
			"message": "etcd get ReplicaSet failed",
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func DeleteReplicaSet(c *gin.Context, s *Server) {
	prefix := "[api-server] [ReplicaSetHandler] [DeleteReplicaSet]"
	fmt.Println(prefix)
	err := c.Request.ParseForm()
	if err != nil {
		return
	}
	if c.Query("all") == "true" {
		// delete the keys
		num, err := s.Etcdstore.DelAll(apiconfig.REPLICASET_PATH)
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

	Name := c.Query("Name")
	fmt.Println("Name:", Name)
	key := c.Request.URL.Path + "/" + Name
	err = s.Etcdstore.Del(key)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delete replicaset failed",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "delete replicaset success",
		"deletePodName": Name,
	})
}
