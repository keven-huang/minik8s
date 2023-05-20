package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
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
		fmt.Println("http get error:", err)
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
	fmt.Println(resList)
	return resList
}

func Watch(resourses string) WatchInterface {
	watcher := &watcher{}
	watcher.resultChan = make(chan Event)
	reader := func(wc chan<- Event) {
		fmt.Println("start watch")
		url := "http://127.0.0.1:8080/watch" + resourses + "?prefix=true"
		resp, err := http.Get(url)
		if err != nil {
			// handle error
			fmt.Println("http get error:", err)
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
				fmt.Println("unmarshal:", event.Key, event.Val, event.Type)
				// send event to watcher.resultChan
				wc <- event
			} else {
				fmt.Println("break")
				break
			}
			time.Sleep(1 * time.Second)
		}
		// This doesn't work(don't know why)
		// reader := bufio.NewReader(resp.Body)
		// for {
		// 	line, err := reader.ReadString('\n')
		// 	if len(line) > 0 {
		// 		fmt.Println("getline")
		// 		// handle Watch Response
		// 		fmt.Println(line)
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
		// 		fmt.Println("break")
		// 		break
		// 	}
		// }
	}
	go reader(watcher.resultChan)
	return watcher
}

func AddPod(pod *core.Pod) error {
	url := apiconfig.Server_URL + apiconfig.POD_PATH
	data, err := json.Marshal(pod)
	if err != nil {
		fmt.Println("failed to marshal person:", err)
		return err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		fmt.Println("add pod success")
	} else {
		return fmt.Errorf("add pod failed")
	}
	return nil
}

func UpdatePod(pod *core.Pod) error {
	url := apiconfig.Server_URL + apiconfig.POD_PATH
	data, err := json.Marshal(pod)
	if err != nil {
		fmt.Println("failed to marshal person:", err)
		return err
	} else {
		fmt.Println("update pod:", string(data))
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()
	fmt.Println("Response Status:", resp.Status)

	return nil
}

func AddNode(node *core.Node) error {
	url := apiconfig.Server_URL + apiconfig.NODE_PATH
	data, err := json.Marshal(node)
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

func GetJobFile(JobName string) core.JobUpload {
	url := apiconfig.Server_URL + apiconfig.JOB_FILE_PATH
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Get Job error", err)
		return core.JobUpload{}
	}
	defer resp.Body.Close()
	reader := resp.Body
	data, err := io.ReadAll(reader)
	if err != nil {
		return core.JobUpload{}
	}
	var job core.JobUpload
	err = json.Unmarshal(data, &job)
	if err != nil {
		return core.JobUpload{}
	}
	return job
}
