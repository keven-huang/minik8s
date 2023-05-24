package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/service"
	"minik8s/pkg/util/log"
	"net/http"
	"time"
)

type EventType string

const (
	Added    EventType = "ADDED"
	Modified EventType = "MODIFIED"
	Deleted  EventType = "DELETED"
	Error    EventType = "ERROR"
)

type Event struct {
	Type EventType
	Key  string
	Val  string
}

type ListRes struct {
	ResourceVersion int64
	Key             string
	Value           string
}

func List(resource string) []ListRes {
	url := "http://127.0.0.1:8080" + resource + "?all=true"
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		fmt.Println("[httpclient] [List] web get error:", err)
	}
	defer resp.Body.Close()
	reader := resp.Body
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil
	}
	var resList []ListRes
	err = json.Unmarshal(data, &resList)
	if err != nil {
		return nil
	}
	fmt.Println("[httpclient] [List] ", resList)
	return resList
}

func Watch(resourses string) WatchInterface {
	watcher := &watcher{}
	watcher.resultChan = make(chan Event)
	reader := func(wc chan<- Event) {
		fmt.Println("[httpclient] [Watch] start watch")
		url := "http://127.0.0.1:8080/watch" + resourses + "?prefix=true"
		resp, err := http.Get(url)
		if err != nil {
			// handle error
			fmt.Println("[httpclient] [Watch] web get error:", err)
		}
		defer resp.Body.Close()
		buf := make([]byte, 40960)
		// Need to optimize
		for {
			n, err := resp.Body.Read(buf)
			if n != 0 || err != io.EOF {
				event := Event{}
				err = json.Unmarshal([]byte(buf[:n]), &event)
				if err != nil {
					continue
				}
				fmt.Println("[httpclient] [Watch] unmarshal:", event.Key, event.Val, event.Type)
				// send event to watcher.resultChan
				wc <- event
			} else {
				fmt.Println("[httpclient] [Watch] break")
				break
			}
			time.Sleep(1 * time.Second)
		}
		// This doesn't work(don't know why)
		// reader := bufio.NewReader(resp.Body)
		// for {
		// 	line, err := reader.ReadString('\n')
		// 	if len(line) > 0 {
		// 		fmt.Println("[httpclient] [Watch] getline")
		// 		// handle Watch Response
		// 		fmt.Println("[httpclient] [Watch] ", line)
		// 		event := Event{}
		// 		// json.Unmarshal([]byte(line), &event)
		// 		// TO DO: send event to watcher.resultChan
		// 		wc <- event
		// 	}
		// 	if err == io.EOF {
		// 		break
		// 	}
		// 	if err != nil {
		// 		// disconnect , cause watch is controlled by client,should try to reconnect
		// 		// TO DO: reconnect
		// 		fmt.Println("[httpclient] [Watch] break")
		// 		break
		// 	}
		// }
	}
	go reader(watcher.resultChan)
	return watcher
}

func UpdatePod(pod *core.Pod) error {
	url := apiconfig.Server_URL + apiconfig.POD_PATH
	data, err := json.Marshal(pod)
	if err != nil {
		fmt.Println("[httpclient] [UpdatePod] failed to marshal person:", err)
		return err
	} else {
		fmt.Println("[httpclient] [UpdatePod] update pod:", string(data))
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("[httpclient] [UpdatePod] Error:", err)
		return err
	}
	defer resp.Body.Close()
	fmt.Println("[httpclient] [UpdatePod] Response Status:", resp.Status)

	return nil
}

func AddNode(node *core.Node) error {
	url := apiconfig.Server_URL + apiconfig.NODE_PATH
	data, err := json.Marshal(node)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	err = log.CheckHttpStatus("[httpclient] [AddNode] ", resp)
	if err != nil {
		return err
	}
	return nil
}
func GetService(name string) (*service.Service, error) {
	prefix := "[tool][GetService]"
	fmt.Println(prefix + "key:" + name)
	url := apiconfig.Server_URL + apiconfig.SERVICE_PATH + "?" + "Name=" + name
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := make([]byte, 40960)
	res := service.Service{}
	n, err := resp.Body.Read(buf)
	if n != 0 || err != io.EOF {
		err = json.Unmarshal([]byte(buf[:n]), &res)
		if err != nil {
			return nil, err
		}
	}
	return &res, nil
}

func UpdateService(service *service.Service) error {
	url := apiconfig.Server_URL + apiconfig.SERVICE_PATH
	fmt.Println("[tool][updateservice]: url=" + url)
	data, err := json.Marshal(service)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("Response Status:", resp.Status)
	return nil
}

func UpdateDNS(dns *core.DNS) error {
	url := apiconfig.Server_URL + apiconfig.DNS_PATH
	fmt.Println("[tool][updateDNS]: url=" + url)
	data, err := json.Marshal(dns)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("Response Status:", resp.Status)
	return nil
}

// TODO 讨论确定一下api-sver的rest-api用法
func DeleteService(service *service.Service) error {
	url := apiconfig.Server_URL + apiconfig.SERVICE_PATH + "/" + service.ServiceSpec.Name
	fmt.Println("[tool][deleteService]: url=" + url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	fmt.Println("Response Status:", resp.Status)
	return nil
}

// Get Pod
// TODO 讨论确定一下具体写法, api路径等

func GetPod(name string) (*core.Pod, error) {
	prefix := "[tool][GetPod]"
	fmt.Println(prefix + "key:" + name)
	url := apiconfig.Server_URL + apiconfig.POD_PATH + name
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := make([]byte, 40960)
	res := core.Pod{}
	n, err := resp.Body.Read(buf)
	if n != 0 || err != io.EOF {
		err = json.Unmarshal([]byte(buf[:n]), &res)
		if err != nil {
			return nil, err
		}
	}
	return &res, nil
}

func GetTypeName(event Event) string {
	var s string
	switch event.Type {
	case Added:
		{
			s = "Add"
		}
	case Modified:
		{
			s = "Modify"
		}
	case Deleted:
		{
			s = "Delete"
		}
	case Error:
		{
			s = "Error"
		}
	}
	return s
}
