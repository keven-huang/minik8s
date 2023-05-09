package kube_service

import (
	"encoding/json"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/service"
	"sync"
	"time"
)

// findPods
// runtimeService根据seletors的信息找到
func (rs *RuntimeService) findPods(isCreate bool) error {
	// TODO
	// 使用informer的list方法获取所有的pods，放到str中
	str = "TODO"
	var curPods []*core.Pod
	// 获取最新的pod消息，放到curPods中,
	for _, val := range str {
		tmp := &core.Pod{}
		err := json.Unmarshal([]byte(val), tmp)
		curPods = append(curPods, tmp)
	}
	// 按照selector筛选pods
	var filtedPods []*core.Pod
	for _, pod := range curPods {
		if pod.Status.Phase != core.PodRunning {
			continue
		}
		isChosed := true
		if rs.ServiceConfig.ServiceSpec.Selector != nil {
			// TODO, pods目前缺少label字段
			for k, v := range rs.ServiceConfig.ServiceSpec.Selector {
				value, ok := pod.Label[k]
				if !ok || v != value{
					isChosed = false
					break
				}
			}
		}
		if isChosed {
			filtedPods = append(filtedPods, pod)
		}
	}
	// 更新pods为最新
	rs.Pods = filtedPods
	// TODO
	// 使用client发送更新的通知给api-server，使其更新etcd
}

// 根据spec创建对应的service，绑定handler
func CreateService(sc *service.Service, informer *informer.Informer) *RuntimeService {
	res := &RuntimeService{}
	res.ServiceConfig = sc
	res.eventChan = make(chan string)
	//TODO, init the informer
	res.Informer = informer
	var lock sync.RWMutex
	res.lock = lock
	res.ifSend = false
	go res.Run(res.eventChan)
	res.findPods(true)
	res.ifSend = true
	res.startTicker()
	return res
}

// 运行这个runtimeService
// 时刻根据pod的真实状态更新缓存
func (rs *RuntimeService) Run(event <-chan string) {
	for {
		select {
		case cmd, ok := <-event:
			if !ok {
				return
			}
			rs.lock.Lock()
			switch cmd {
			case TICK_EVENT:
				flag := false
				for _, pod := range rs.Pods {
					// 应该调用restClient，根据pod的名字查询当前pod的状态
					// TODO
					cur_pod, err := new core.Pod{}, nil
					if err != nil {
						continue
					}
					if cur_pod == nil {
						pod.Status.Phase = core.PodSucceeded
						flag = true
					}
				}
				if (flag) {
					err := rs.findPods(false)
				}
				rs.ifSend = true
			}
			rs.lock.Unlock()
		}
	}
}

// 启动一个协程，每10秒发送一个tick, tick会被run捕获从而更新RunService中的状态
func (rs *RuntimeService)startTicker()  {
	rs.timer = time.NewTicker(10 * time.Second)
	go func(rs *RuntimeService) {
		defer rs.timer.Stop()
		for {
			select {
			case <- rs.timer.C:
				rs.lock.Lock()
				if rs.ifSend {
					rs.eventChan <- TICK_EVENT
					rs.ifSend = false
				}
				rs.lock.Unlock()
			}
		}
	}(rs)
}

// delete current runtimeService
func (*RuntimeService) Delete() {

}
