package server

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Put(c *gin.Context, s *Server) {
	fmt.Printf("put\n")
	key := c.Request.URL.Path
	value, _ := io.ReadAll(c.Request.Body)
	err := s.Etcdstore.Put(key, string(value))
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "put success",
	})
}

func Get(c *gin.Context, s *Server) {
	key := c.Request.URL.Path
	res, err := s.Etcdstore.Get(key)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": res,
	})
}

func Delete(c *gin.Context, s *Server) {
	key := c.Request.URL.Path
	err := s.Etcdstore.Del(key)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "delete success",
	})
}

func Post(c *gin.Context, s *Server) {
	key := c.Request.URL.Path
	value, _ := io.ReadAll(c.Request.Body)
	err := s.Etcdstore.Put(key, string(value))
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "post success",
	})
}
