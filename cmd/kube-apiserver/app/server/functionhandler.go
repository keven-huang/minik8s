package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"minik8s/pkg/api/core"
	v1 "minik8s/pkg/apis/meta/v1"
	"minik8s/pkg/util/random"
	"net/http"
)

func GetFunction(c *gin.Context, s *Server) {

}

func DeleteFunction(c *gin.Context, s *Server) {

}

func AddFunction(c *gin.Context, s *Server) {
	prefix := "[api-server] [AddFunction] "
	val, _ := io.ReadAll(c.Request.Body)
	function := core.Function{}
	err := json.Unmarshal([]byte(val), &function)
	if err != nil {
		log.Println("[ERROR] ", prefix, err)
		return
	}
	key := c.Request.URL.Path + "/" + function.Name
	res, _ := s.Etcdstore.Get(key)
	if len(res) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "function Name Duplicate.",
		})
		return
	}

	function.UID = random.GenerateUUID()
	function.ObjectMeta.CreationTimestamp = v1.Now()

	body, _ := json.Marshal(function)

	err = s.Etcdstore.Put(key, string(body))
	if err != nil {
		log.Println("[ERROR] ", prefix, err)
		return
	}
}
