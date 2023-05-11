package kubelet

import (
	"encoding/json"
	"fmt"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	"minik8s/pkg/kubelet/dockerClient"
)

type Kublet struct {
	PodInformer informer.Informer
}

func NewKublet() *Kublet {
	return &Kublet{
		PodInformer: informer.NewInformer(apiconfig.POD_PATH),
	}
}

func (k *Kublet) Register() {
	k.PodInformer.AddEventHandler(tool.Added, k.CreatePod)
	k.PodInformer.AddEventHandler(tool.Deleted, k.DeletePod)
}

func (k *Kublet) CreatePod(event tool.Event) {
	// handle event
	fmt.Println("In AddPod EventHandler:")
	fmt.Println("event.Key: ", event.Key)
	fmt.Println("event.Val: ", event.Val)
	k.PodInformer.Set(event.Key, event.Val)

	pod := &core.Pod{}
	err := json.Unmarshal([]byte(event.Val), pod)
	if err != nil {
		fmt.Println(err)
		return
	}

	metaData, netSetting, err := dockerClient.CreatePod(*pod)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, meta := range metaData {
		fmt.Println(meta.Name, meta.Id)
	}
	net, err := json.MarshalIndent(netSetting, "", "  ")
	fmt.Println("net: ")
	fmt.Println(string(net))
	fmt.Println("-----------")
}

func (k *Kublet) DeletePod(event tool.Event) {
	// handle event
	fmt.Println("In DeletePod EventHandler:")
	fmt.Println("event.Key: ", event.Key)
	fmt.Println("event.Val: ", k.PodInformer.Get(event.Key))

	pod := &core.Pod{}
	err := json.Unmarshal([]byte(k.PodInformer.Get(event.Key)), pod)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = dockerClient.DeletePod(*pod)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (k *Kublet) Run() {
	k.PodInformer.Run()
}
