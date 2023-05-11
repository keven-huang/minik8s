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

type Kubelet struct {
	PodInformer informer.Informer
	node        core.Node
}

func NewKubelet(name string) (*Kubelet, error) {
	node := core.Node{}
	node.Name = name
	err := tool.AddNode(&node)
	if err != nil {
		return nil, err
	}
	return &Kubelet{
		PodInformer: informer.NewInformer(apiconfig.POD_PATH),
		node:        node,
	}, nil
}

func (k *Kubelet) Register() {
	k.PodInformer.AddEventHandler(tool.Added, k.CreatePod)
	k.PodInformer.AddEventHandler(tool.Modified, k.CreatePod)
	k.PodInformer.AddEventHandler(tool.Deleted, k.DeletePod)
}

func (k *Kubelet) CreatePod(event tool.Event) {
	// handle event
	fmt.Println("In AddPod/ModifyPod EventHandler: ", tool.GetTypeName(event))
	fmt.Println("event.Key: ", event.Key)
	fmt.Println("event.Val: ", event.Val)
	k.PodInformer.Set(event.Key, event.Val)

	pod := &core.Pod{}
	err := json.Unmarshal([]byte(event.Val), pod)
	if err != nil {
		fmt.Println(err)
		return
	}

	if pod.Spec.NodeName != k.node.Name {
		fmt.Println("Not mine. Node:", pod.Spec.NodeName)
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

func (k *Kubelet) DeletePod(event tool.Event) {
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

func (k *Kubelet) Run() {
	k.PodInformer.Run()
}
