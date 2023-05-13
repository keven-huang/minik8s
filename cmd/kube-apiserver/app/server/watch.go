package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Watch(c *gin.Context, s *Server) {
	key := c.Request.URL.Path[6:]
	fmt.Println("path:", c.Request.URL.Path)
	fmt.Println("key:", key)

	var isPrefix = true
	if c.Query("prefix") == "false" {
		isPrefix = false
	}

	resChan := s.Etcdstore.Watch(key, isPrefix)
	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Printf("http server does not support flush\n")
		return
	}
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	flusher.Flush()
	fmt.Println("in sender")
	for {
		select {
		case res, ok := <-resChan:
			if !ok {
				// resChan 已关闭，退出循环
				return
			}
			fmt.Println("send watch response")
			resp, _ := json.Marshal(res)
			fmt.Println(string(resp))
			fmt.Fprintf(w, string(resp))
			flusher.Flush()
		}
	}
}
