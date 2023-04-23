package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Put(c *gin.Context) {
	fmt.Printf("put\n")
	key := c.PostForm("key")
	value := c.PostForm("value")
	fmt.Printf("key: %s,value: %s", key, value)
	err := s.etcdstore.Put(key, value)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "put success",
	})
}

func (s *Server) Get(c *gin.Context) {
	key := c.PostForm("key")
	fmt.Printf("key: %s", key)
	res, err := s.etcdstore.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
}

func (s *Server) Delete(c *gin.Context) {
	key := c.PostForm("key")
	fmt.Printf("key: %s", key)
	err := s.etcdstore.Del(key)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "delete success",
	})
}

func (s *Server) Update(c *gin.Context) {

}
