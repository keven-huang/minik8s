package kubelet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	"minik8s/pkg/kubelet/dockerClient"
	"minik8s/pkg/util/web"
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
	prefix := "[kubelet] [CreatePod] "
	// handle event
	fmt.Println("In AddPod/ModifyPod EventHandler: ", tool.GetTypeName(event))
	fmt.Println("event.Key: ", event.Key)
	fmt.Println("event.Val: ", event.Val)
	k.PodInformer.Set(event.Key, event.Val)

	pod := &core.Pod{}
	err := json.Unmarshal([]byte(event.Val), pod)
	if err != nil {
		fmt.Println(prefix, err)
		return
	}

	if pod.Spec.NodeName != k.node.Name {
		fmt.Println(prefix, "Not mine. The node of the pod is :", pod.Spec.NodeName)
		return
	}

	if pod.Status.Phase != "Pending" {
		fmt.Println(prefix, "phase is not satisfied:", pod.Status.Phase)
		return
	}

	for i, v := range pod.Spec.Containers {
		pod.Spec.Containers[i].Name = pod.Name + "-" + v.Name
	}

	metaData, netSetting, err := dockerClient.CreatePod(*pod)
	if err != nil {
		fmt.Println(err)
		return
	}

	pod.Status.Phase = "Running"
	// 创建成功 修改Status
	data, err := json.Marshal(pod)
	if err != nil {
		fmt.Println(prefix, "failed to marshal:", err)
	}

	err = web.SendHttpRequest("POST", apiconfig.Server_URL+apiconfig.POD_PATH,
		web.WithPrefix(prefix), web.WithBody(bytes.NewBuffer(data)))
	if err != nil {
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
	k.PodInformer.Delete(event.Key)

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
