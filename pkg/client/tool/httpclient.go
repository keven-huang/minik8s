package tool

import (
	"encoding/json"
	"fmt"
	"io"
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
	ResourceVersion string
	Key             string
	Value           string
}

func List(resource string) []ListRes {
	url := "http://127.0.0.1:8080" + resource
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
	return resList
}

func Watch(resourses string) WatchInterface {
	watcher := &watcher{}
	watcher.resultChan = make(chan Event)
	reader := func(wc chan<- Event) {
		fmt.Println("start watch")
		//url := "http://127.0.0.1:8080" + resourses + "?watch=true"
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
