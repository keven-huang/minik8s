package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"minik8s/pkg/api/core"
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
