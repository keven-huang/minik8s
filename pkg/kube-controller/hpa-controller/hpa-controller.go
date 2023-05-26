package hpacontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	"minik8s/pkg/kubelet/monitor"
	q "minik8s/pkg/util/concurrentqueue"
	"minik8s/pkg/util/web"
	"net/http"
	"time"
)

type HPAController struct {
	PodInformer        informer.Informer
	ReplicasetInformer informer.Informer
	HPAInformer        informer.Informer
	queue              *q.ConcurrentQueue
	mark               map[string]bool
}

func NewHPAController() *HPAController {
	return &HPAController{
		PodInformer:        informer.NewInformer(apiconfig.POD_PATH),
		ReplicasetInformer: informer.NewInformer(apiconfig.REPLICASET_PATH),
		HPAInformer:        informer.NewInformer(apiconfig.HPA_PATH),
		queue:              q.NewConcurrentQueue(),
		mark:               make(map[string]bool),
	}
}

func (hpac *HPAController) Register() {
	hpac.PodInformer.AddEventHandler(tool.Added, hpac.AddPod)
	hpac.PodInformer.AddEventHandler(tool.Modified, hpac.UpdatePod)
	hpac.PodInformer.AddEventHandler(tool.Deleted, hpac.DeletePod)
	hpac.ReplicasetInformer.AddEventHandler(tool.Added, hpac.AddReplicaset)
	hpac.ReplicasetInformer.AddEventHandler(tool.Modified, hpac.UpdateReplicaset)
	hpac.ReplicasetInformer.AddEventHandler(tool.Deleted, hpac.DeleteReplicaset)
	hpac.HPAInformer.AddEventHandler(tool.Added, hpac.AddHPA)
	//hpac.HPAInformer.AddEventHandler(tool.Modified, hpac.UpdateHPA)
	hpac.HPAInformer.AddEventHandler(tool.Deleted, hpac.DeleteHPA)
}

func (hpac *HPAController) Run() {
	go hpac.PodInformer.Run()
	go hpac.ReplicasetInformer.Run()
	go hpac.HPAInformer.Run()
	go worker(hpac)
	select {}
}

func updateReplicaSetToServer(hpac *HPAController, replica *core.ReplicaSet, key string, prefix string) error {
	data, err := json.Marshal(replica)
	if err != nil {
		fmt.Println(prefix, "failed to marshal:", err)
		return err
	}

	err = web.SendHttpRequest("POST", apiconfig.Server_URL+apiconfig.REPLICASET_PATH,
		web.WithPrefix(prefix), web.WithBody(bytes.NewBuffer(data)))
	if err != nil {
		return err
	}

	hpac.ReplicasetInformer.Set(key, string(data))
	return nil
}

func Scale(hpac *HPAController, r *core.ReplicaSet, r_key string, from int32, to int32, t int32) error {
	fmt.Println("[HPA] [Scale] from: ", from, " to: ", to, " t: ", t, " r.Name: ", r.Name)
	if from < to {
		for i := from + 1; i <= to; i++ {
			*r.Spec.Replicas = i
			//更新ReplicaSet
			err := updateReplicaSetToServer(hpac, r, r_key, "[HPA] [Scale] [updateReplicaSetToServer]")
			if err != nil {
				fmt.Println("[HPA] [Scale] [updateReplicaSetToServer] err: ", err)
				return err
			}
			time.Sleep(time.Duration(t) * time.Second)
		}
	} else if from > to {
		for i := from - 1; i >= to; i-- {
			*r.Spec.Replicas = i
			//更新ReplicaSet
			err := updateReplicaSetToServer(hpac, r, r_key, "[HPA] [Scale] [updateReplicaSetToServer]")
			if err != nil {
				fmt.Println("[HPA] [Scale] [updateReplicaSetToServer] err: ", err)
				return err
			}
			time.Sleep(time.Duration(t) * time.Second)
		}
	}

	return nil
}

func calcNewReplicas(hpac *HPAController, hpa *core.HPA, rs *core.ReplicaSet) int {
	res := tool.List(apiconfig.NODE_PATH)

	results := make([]monitor.StatsResponse, 0)

	for _, item := range res {
		node := &core.Node{}
		err := json.Unmarshal([]byte(item.Value), node)
		if err != nil {
			fmt.Println("[HPA] [calcNewReplicas] [Unmarshal] err: ", err)
			continue
		}
		url := "http://" + node.Spec.NodeIP + ":9400/stats"
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("[HPA] [calcNewReplicas] [http.Get] err: ", err)
			continue
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("[HPA] [calcNewReplicas] [ioutil.ReadAll] err: ", err)
			continue
		}
		responds := make([]monitor.StatsResponse, 0)
		err = json.Unmarshal(body, &responds)
		if err != nil {
			fmt.Println("[HPA] [calcNewReplicas] [Unmarshal StatsResponse] err: ", err)
			continue
		}

		for _, respond := range responds {
			pod_name := respond.PodName
			pod := &core.Pod{}
			pod_key := apiconfig.POD_PATH + "/" + pod_name
			err = json.Unmarshal([]byte(hpac.PodInformer.Get(pod_key)), pod)
			if err != nil {
				fmt.Println("[HPA] [calcNewReplicas] [Unmarshal Pod] err: ", err)
				continue
			}
			for _, owner := range pod.OwnerReferences {
				if owner.Kind == "ReplicaSet" && owner.Name == rs.Name {
					results = append(results, respond)
					fmt.Println("[HPA] [calcNewReplicas] [append] nodeName: ", node.Name, " podName: ", pod.Name)
					break
				}
			}
		}
	}

	if len(results) == 0 {
		fmt.Println("[HPA] [calcNewReplicas] [len(results) == 0]")
		return int(hpa.Spec.MinReplicas)
	}

	//计算平均值
	cpu_sum := 0.0
	mem_sum := 0.0
	for _, result := range results {
		cpu_sum += result.CPUUtilization
		mem_sum += result.MemoryUsage
		fmt.Println("[HPA] [calcNewReplicas] podName: ", result.PodName, " cpu: ", result.CPUUtilization, " mem: ", result.MemoryUsage)
	}
	fmt.Println("[HPA] [calcNewReplicas] cpu_sum: ", cpu_sum, " mem_sum: ", mem_sum)

	num := 0

	for _, metric := range hpa.Spec.Metrics {
		var target int
		if metric.Name == "cpu" {
			target = int(math.Ceil(cpu_sum / float64(metric.Target.Value)))
		}
		if metric.Name == "memory" {
			target = int(math.Ceil(mem_sum / float64(metric.Target.Value)))
		}
		fmt.Println("[HPA] [calcNewReplicas] target: ", target, " metric.Name: ", metric.Name)
		// conpare target with minReplicas and maxReplicas
		if target < int(hpa.Spec.MinReplicas) {
			target = int(hpa.Spec.MinReplicas)
		}
		if target > int(hpa.Spec.MaxReplicas) {
			target = int(hpa.Spec.MaxReplicas)
		}
		if num < target {
			num = target
		}
	}
	fmt.Println("[HPA] [calcNewReplicas] NewReplicas: ", num)
	return num
}

func worker(hpac *HPAController) {
	periodTime := 20 * time.Second
	for {
		for _, v := range *hpac.HPAInformer.GetCache() {
			hpa := &core.HPA{}
			err := json.Unmarshal([]byte(v), hpa)
			if err != nil {
				fmt.Println("[HPA] [worker] [Unmarshal] err: ", err)
				continue
			}

			if hpac.mark[hpa.Name] {
				fmt.Println("[HPA] [worker] [on mark] hpa.Name: ", hpa.Name)
				continue
			}

			r := &core.ReplicaSet{}
			r_key := apiconfig.REPLICASET_PATH + "/" + hpa.Spec.ScaleTargetRef.Name
			err = json.Unmarshal([]byte(hpac.ReplicasetInformer.Get(r_key)), r)
			if err != nil {
				fmt.Println("[HPA] [worker] [Unmarshal] err: ", err)
				continue
			}

			//获取属于这个HPA对应的ReplicaSet的所有Pod的对应指标的平均数=
			// 根据 hpa.Spec.ScaleTargetRef.Name 和 Pod的OwnerReferences判断
			newReplicas := calcNewReplicas(hpac, hpa, r)

			go func() {
				hpac.mark[hpa.Name] = true
				err := Scale(hpac, r, r_key, r.Status.Replicas, int32(newReplicas), hpa.Spec.PeriodSeconds)
				if err != nil {
					fmt.Println("[HPA] [worker] [Scale] err: ", err)
				}
				hpac.mark[hpa.Name] = false
			}()
		}

		time.Sleep(periodTime)
	}
}

func (hpac *HPAController) AddPod(event tool.Event) {
	prefix := "[HPA] [AddPod] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	hpac.PodInformer.Set(event.Key, event.Val)
}

func (hpac *HPAController) UpdatePod(event tool.Event) {
	prefix := "[HPA] [UpdatePod] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	hpac.PodInformer.Set(event.Key, event.Val)
}

func (hpac *HPAController) DeletePod(event tool.Event) {
	prefix := "[HPA] [DeletePod] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	//val := hpac.PodInformer.Get(event.Key)
	hpac.PodInformer.Delete(event.Key)
}

func (hpac *HPAController) AddReplicaset(event tool.Event) {
	prefix := "[HPA] [AddReplicaset] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	hpac.ReplicasetInformer.Set(event.Key, event.Val)
}

func (hpac *HPAController) UpdateReplicaset(event tool.Event) {
	prefix := "[HPA] [UpdateReplicaset] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	hpac.ReplicasetInformer.Set(event.Key, event.Val)
}

func (hpac *HPAController) DeleteReplicaset(event tool.Event) {
	prefix := "[HPA] [DeleteReplicaset] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	//val := hpac.ReplicasetInformer.Get(event.Key)
	hpac.ReplicasetInformer.Delete(event.Key)
}

func (hpac *HPAController) AddHPA(event tool.Event) {
	prefix := "[HPA] [AddHPA] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	hpac.HPAInformer.Set(event.Key, event.Val)

	hpa := &core.HPA{}
	err := json.Unmarshal([]byte(event.Val), hpa)
	if err != nil {
		fmt.Println("[HPA] [AddHPA] [Unmarshal] err: ", err)
		return
	}
	hpac.mark[hpa.Name] = false
}

func (hpac *HPAController) DeleteHPA(event tool.Event) {
	prefix := "[HPA] [DeleteHPA] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	//val := hpac.HPAInformer.Get(event.Key)
	hpac.HPAInformer.Delete(event.Key)
}
