package scheduler

import (
	"encoding/json"
	"fmt"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	q "minik8s/pkg/util/concurrentqueue"
	"time"
)

type Scheduler struct {
	PodInformer  informer.Informer
	NodeInformer informer.Informer
	queue        *q.ConcurrentQueue
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		PodInformer:  informer.NewInformer(apiconfig.POD_PATH),
		NodeInformer: informer.NewInformer(apiconfig.NODE_PATH),
		queue:        q.NewConcurrentQueue(),
	}
}

func (s *Scheduler) Register() {
	s.PodInformer.AddEventHandler(tool.Added, s.AddPod)
}

func (s *Scheduler) AddPod(event tool.Event) {
	fmt.Println("add pod")
	pod := &core.Pod{}
	err := json.Unmarshal([]byte(event.Val), &pod)
	if err != nil {
		return
	}
	if pod.Spec.NodeName != "" {
		return
	}
	s.queue.Push(pod)
}

func (s *Scheduler) worker() {
	for {
		//TODO: optmize can use channel or condition variable
		if !s.queue.IsEmpty() {
			pod := s.queue.Pop()
			s.Schedule(pod.(*core.Pod))
		} else {
			time.Sleep(time.Second)
		}
	}
}

func (s *Scheduler) Schedule(pod *core.Pod) {
	node := s.GetNode()
	var nodeName string
	nodeName = roundrobin_strategy(node)
	pod.Spec.NodeName = nodeName
	fmt.Println("schedule to node:", nodeName)
	tool.UpdatePod(pod)
}

func (s *Scheduler) GetNode() []core.Node {
	// TODO : can optimize by cache
	res := s.NodeInformer.GetCache()
	var nodes []core.Node
	for _, val := range *res {
		node := core.Node{}
		err := json.Unmarshal([]byte(val), &node)
		if err != nil {
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func (s *Scheduler) Run() {
	go s.PodInformer.Run()
	go s.worker()
	select {}
}
