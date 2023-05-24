package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	"minik8s/pkg/util/web"
	"time"
)

func CreatePod(event tool.Event) {
	fmt.Println("[test] [CreatePod] event.Type: ", tool.GetTypeName(event))
	time.Sleep(3 * time.Second)
}

func send() {
	time.Sleep(3 * time.Second)
	for i := 0; i < 3; i++ {
		pod := &core.Pod{}
		pod.Name = "test" + fmt.Sprint(i)

		data, err := json.Marshal(pod)
		if err != nil {
			fmt.Println(err)
		}
		err = web.SendHttpRequest("PUT", apiconfig.Server_URL+apiconfig.POD_PATH,
			web.WithPrefix("[test] [send] "),
			web.WithBody(bytes.NewBuffer(data)),
			web.WithLog(true))
		if err != nil {
			return
		}
		// time.Sleep(3 * time.Second)
	}

}

func main() {
	PodInformer := informer.NewInformer(apiconfig.POD_PATH)
	PodInformer.AddEventHandler(tool.Added, CreatePod)

	go send()

	PodInformer.Run()
}
