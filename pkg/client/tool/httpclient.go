package tool

import (
	"encoding/json"
	"fmt"
	"log"
	"minik8s/pkg/kube-apiserver/etcd"
	"net/http"
	"time"
)

type Event etcd.Event

func Watch(resourses string) WatchInterface {
	buf := make([]byte, 1024)
	watcher := &watcher{}
	watcher.resultChan = make(chan Event)
	reader := func(wc chan<- Event) {
		url := "http://127.0.0.1:8080/" + resourses + "?watch=true"
		resp, err := http.Get(url)
		if err != nil {
			// handle error
			log.Println("http get error:", err)
		}
		defer resp.Body.Close()
		for {
			n, err := resp.Body.Read(buf)
			if n != 0 || err != nil {
				// handle Watch Response
				fmt.Println(string(buf[:n]))
				event := Event{}
				json.Unmarshal(buf[:n], &event)
				// TO DO: send event to watcher.resultChan
				wc <- event
			} else {
				// disconnect , cause watch is controlled by client,should try to reconnect
				// TO DO: reconnect
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
	go reader(watcher.resultChan)
	return watcher
}

func List() {

}
