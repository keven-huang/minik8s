package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/kube-apiserver/etcd"
	"net/http"
)

func UpdateDNS(c *gin.Context, s *Server) {
	prefix := "[DNSHandler][UpdateDNS]:"
	val, _ := io.ReadAll(c.Request.Body)
	dns := core.DNS{}
	err := json.Unmarshal([]byte(val), &dns)
	if err != nil {
		fmt.Println(err)
		return
	}
	key := c.Request.URL.Path + "/" + dns.Metadata.Name
	body, _ := json.Marshal(dns)
	fmt.Println(prefix + "key:" + key)
	fmt.Println(prefix + "value:" + string(body))
	err = s.Etcdstore.Put(key, string(body))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "etcd put dns failed.",
		})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "update dns success.",
	})
}

func GetDNS(c *gin.Context, s *Server) {
	fmt.Println("[DnsHandler] [GetDNS]")
	if c.Query("all") == "true" {
		fmt.Println("[warning] should be only one dns config")
		res, err := s.Etcdstore.GetWithPrefix(apiconfig.DNS_PATH)
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}

	DnsName := c.Query("Name")
	key := c.Request.URL.Path + "/" + string(DnsName)

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
			"message": "etcd get dns failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get dns successfully.",
		"Results": res,
	})
}

func DeleteDNS(c *gin.Context, s *Server) {
	prefix := "[DNSHandler][DeleteDNS]:"
	if c.Query("all") == "true" { // delete all dns
		num, err := s.Etcdstore.DelAll(apiconfig.DNS_PATH)
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":   "delete all dns successfully.",
			"deleteNum": num,
		})
		return
	}
	DnsConfigName := c.Query("DnsConfigName")
	key := c.Request.URL.Path + "/" + DnsConfigName
	fmt.Println(prefix + "key:" + key)
	err := s.Etcdstore.Del(key)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delete dns failed",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":           "delete service success",
		"deleteServiceName": DnsConfigName,
	})
}
