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
	"minik8s/pkg/cmd/create"
	"minik8s/pkg/util/random"
	"net/http"
	"strings"
	"time"
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
		// 删除所有之前这个Function版本的Pod
		deletePodWithFunctionName(s, function.Name)
	}

	function.UID = random.GenerateUUID()
	function.ObjectMeta.CreationTimestamp = v1.Now()

	body, _ := json.Marshal(function)

	err = s.Etcdstore.Put(key, string(body))
	if err != nil {
		log.Println("[ERROR] ", prefix, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "add function success.",
	})
}

func deletePodWithFunctionName(s *Server, func_name string) {
	res, err := s.Etcdstore.GetWithPrefix(apiconfig.POD_PATH)
	if err != nil {
		fmt.Println("[ERROR] [deletePodWithFunctionName] ", err)
		return
	}
	for _, v := range res {
		pod := core.Pod{}
		_ = json.Unmarshal([]byte(v.Value), &pod)
		if pod.OwnerReferences[0].Name == func_name && pod.OwnerReferences[0].Kind == "Function" {
			_ = s.Etcdstore.Del(v.Key)
		}
	}
}

func scheduler(s *Server, func_name string) (*core.Pod, error) {
	var p []*core.Pod = make([]*core.Pod, 0)
	res, err := s.Etcdstore.GetWithPrefix(apiconfig.POD_PATH)
	if err != nil {
		fmt.Println("[ERROR] [scheduler] ", err)
		return nil, err
	}
	for _, v := range res {
		pod := &core.Pod{}
		_ = json.Unmarshal([]byte(v.Value), pod)
		if pod.OwnerReferences[0].Name == func_name && pod.OwnerReferences[0].Kind == "Function" {
			p = append(p, pod)
		}
	}
	function := core.Function{}
	r, err := s.Etcdstore.GetExact(apiconfig.FUNCTION_PATH + "/" + func_name)
	if err != nil {
		fmt.Println("[ERROR] [scheduler] [GetExact]", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(r[0].Value), &function)
	if err != nil {
		fmt.Println("[ERROR] [scheduler] [Unmarshal]", err)
		return nil, err
	}
	function.Spec.InvokeTimes++
	data, err := json.Marshal(function)
	if err != nil {
		fmt.Println("[ERROR] [scheduler] [Marshal]", err)
		return nil, err
	}
	err = s.Etcdstore.Put(apiconfig.FUNCTION_PATH+"/"+func_name, string(data))
	if err != nil {
		fmt.Println("[ERROR] [scheduler] [Put]", err)
		return nil, err
	}
	return p[function.Spec.InvokeTimes%len(p)], nil
}

func InvokeFunction(c *gin.Context, s *Server) {
	function_name := c.Param("function_name")[1:]
	pod_name := "function-" + function_name + "-" + random.GenerateRandomString(5)

	pod, err := scheduler(s, function_name)
	var podIP string

	if pod == nil {
		pod = &core.Pod{
			ObjectMeta: v1.ObjectMeta{
				Name: pod_name,
				OwnerReferences: []v1.OwnerReference{
					{
						Kind: "Function",
						Name: function_name,
					},
				},
			},
			Spec: core.PodSpec{
				Containers: []core.Container{
					{
						Name:  pod_name + "-container",
						Image: "luhaoqi/my_module:" + function_name,
					},
				},
			},
		}

		err = create.CreatePod(pod)
		if err != nil {
			fmt.Println("[InvokeFunction] ", "Error creating pod:", err)
			return
		}

		for {
			r, err := s.Etcdstore.GetExact(apiconfig.POD_PATH + "/" + pod.Name)
			if err != nil {
				fmt.Println("[InvokeFunction] ", "Error getting pod:", err)
				return
			}
			if len(r) > 0 {
				if strings.Contains(r[0].Value, `"podIP":"`) {
					break
				}
			}
			time.Sleep(1 * time.Second)
		}

	} else {
		podIP = pod.Status.PodIP
	}

	req, err := http.NewRequest("POST", "http://"+podIP+":8888", c.Request.Body)
	if err != nil {
		fmt.Println("[InvokeFunction] ", "Error creating request:", err)
		return
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[InvokeFunction]", "Error sending request:", err)
		return
	}

	c.JSON(http.StatusOK, resp.Body)
}
