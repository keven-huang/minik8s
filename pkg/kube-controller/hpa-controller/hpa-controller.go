package hpacontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	q "minik8s/pkg/util/concurrentqueue"
	"minik8s/pkg/util/web"
	"time"
)

type HPAController struct {
	PodInformer        informer.Informer
	ReplicasetInformer informer.Informer
	HPAInformer        informer.Informer
	queue              *q.ConcurrentQueue
}

func NewHPAController() *HPAController {
	return &HPAController{
		PodInformer:        informer.NewInformer(apiconfig.POD_PATH),
		ReplicasetInformer: informer.NewInformer(apiconfig.REPLICASET_PATH),
		HPAInformer:        informer.NewInformer(apiconfig.HPA_PATH),
		queue:              q.NewConcurrentQueue(),
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
			r := &core.ReplicaSet{}
			r_key := apiconfig.REPLICASET_PATH + "/" + hpa.Spec.ScaleTargetRef.Name
			err = json.Unmarshal([]byte(hpac.ReplicasetInformer.Get(r_key)), r)
			if err != nil {
				fmt.Println("[HPA] [worker] [Unmarshal] err: ", err)
				continue
			}

			//获取属于这个HPA对应的ReplicaSet的所有Pod的对应指标的平均数=
			// 根据 hpa.Spec.ScaleTargetRef.Name 和 Pod的OwnerReferences判断
			newReplicas := 3

			Scale(hpac, r, r_key, r.Status.Replicas, int32(newReplicas), hpa.Spec.PeriodSeconds)

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
}

func (hpac *HPAController) DeleteHPA(event tool.Event) {
	prefix := "[HPA] [DeleteHPA] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	//val := hpac.HPAInformer.Get(event.Key)
	hpac.HPAInformer.Delete(event.Key)
}
