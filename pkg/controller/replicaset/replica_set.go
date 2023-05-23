package replicaset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/apis/meta/v1"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	"minik8s/pkg/cmd/create"
	q "minik8s/pkg/util/concurrentqueue"
	"minik8s/pkg/util/random"
	"minik8s/pkg/util/web"
	"net/url"
	"time"
)

type ReplicaSetController struct {
	PodInformer        informer.Informer
	ReplicasetInformer informer.Informer
	queue              *q.ConcurrentQueue
}

func NewReplicaSetController() (*ReplicaSetController, error) {
	return &ReplicaSetController{
		PodInformer:        informer.NewInformer(apiconfig.POD_PATH),
		ReplicasetInformer: informer.NewInformer(apiconfig.REPLICASET_PATH),
		queue:              q.NewConcurrentQueue(),
	}, nil
}

func (rsc *ReplicaSetController) Register() {
	rsc.PodInformer.AddEventHandler(tool.Added, rsc.AddPod)
	rsc.PodInformer.AddEventHandler(tool.Modified, rsc.UpdatePod)
	rsc.PodInformer.AddEventHandler(tool.Deleted, rsc.DeletePod)
	rsc.ReplicasetInformer.AddEventHandler(tool.Added, rsc.AddReplicaset)
	rsc.ReplicasetInformer.AddEventHandler(tool.Modified, rsc.UpdateReplicaset)
	rsc.ReplicasetInformer.AddEventHandler(tool.Deleted, rsc.DeleteReplicaset)
}

// Match 判断Pod是否符合ReplicaSet的条件
func Match(rs *core.ReplicaSet, pod *core.Pod) bool {
	// A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
	for key, value := range rs.Spec.Selector.MatchLabels {
		v, exists := pod.Labels[key]
		if !exists || v != value {
			return false
		}
	}
	for _, expression := range rs.Spec.Selector.MatchExpressions {
		// Valid operators are In, NotIn, Exists and DoesNotExist.
		key := expression.Key
		val, exists := pod.Labels[key]
		flag := false
		for _, v := range expression.Values {
			if v == val {
				flag = true
				break
			}
		}
		switch expression.Operator {
		case v1.LabelSelectorOpIn:
			{
				if !flag {
					return false
				}
			}
		case v1.LabelSelectorOpNotIn:
			{
				if flag {
					return false
				}
			}
		case v1.LabelSelectorOpExists:
			{
				if !exists {
					return false
				}
			}
		case v1.LabelSelectorOpDoesNotExist:
			{
				if exists {
					return false
				}
			}
		default:
			{
				fmt.Println("[ERROR] [ReplicaSet] [Match] Operator not supported: ", expression.Operator)
			}
		}
	}
	return true
}

func updateReplicaSet(rsc *ReplicaSetController, replica *core.ReplicaSet, key string, prefix string) {
	data, err := json.Marshal(replica)
	if err != nil {
		fmt.Println(prefix, "failed to marshal:", err)
	}

	err = web.SendHttpRequest("POST", apiconfig.Server_URL+apiconfig.REPLICASET_PATH,
		web.WithPrefix(prefix), web.WithBody(bytes.NewBuffer(data)))
	if err != nil {
		return
	}

	rsc.ReplicasetInformer.Set(key, string(data))
}

func (rsc *ReplicaSetController) AddReplicaset(event tool.Event) {
	prefix := "[ReplicaSet] [AddReplicaset] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	fmt.Println(prefix, "event.Val: ", event.Val)
	rsc.ReplicasetInformer.Set(event.Key, event.Val)

	replica := &core.ReplicaSet{}
	err := json.Unmarshal([]byte(event.Val), replica)
	if err != nil {
		fmt.Println("[ERROR] ", prefix, err)
		return
	}

	if replica.Status.Replicas != 0 {
		fmt.Println("[ERROR] ", prefix, "replica.Status.Replicas: ", replica.Status.Replicas)
		return
	}

	pod_cache := rsc.PodInformer.GetCache()
	flag := false

	for _, value := range *pod_cache {
		pod := &core.Pod{}
		err := json.Unmarshal([]byte(value), pod)
		if err != nil {
			log.Println("[ERROR] ", prefix, err)
			return
		}
		if (pod.Status.Phase == "Running") && Match(replica, pod) {
			pod.OwnerReferences = append(pod.OwnerReferences,
				v1.OwnerReference{
					Name: replica.Name,
					UID:  replica.UID,
				})
			replica.Status.Replicas++
			flag = true
			fmt.Println(prefix, "find pod: ", pod.Name)
		}
	}

	if flag {
		updateReplicaSet(rsc, replica, event.Key, prefix)
	}

	if *replica.Spec.Replicas != replica.Status.Replicas {
		rsc.queue.Push(event.Key)
	}

}

func (rsc *ReplicaSetController) UpdateReplicaset(event tool.Event) {
	prefix := "[ReplicaSet] [UpdateReplicaset] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	fmt.Println(prefix, "event.Val: ", event.Val)
	rsc.ReplicasetInformer.Set(event.Key, event.Val)
}

func (rsc *ReplicaSetController) DeleteReplicaset(event tool.Event) {
	prefix := "[ReplicaSet] [DeleteReplicaset] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	fmt.Println(prefix, "event.Val: ", event.Val)

	val := rsc.ReplicasetInformer.Get(event.Key)
	rsc.ReplicasetInformer.Delete(event.Key)

	fmt.Println(prefix, "val: ", val)
	replica := &core.ReplicaSet{}
	err := json.Unmarshal([]byte(val), replica)
	if err != nil {
		fmt.Println("[ERROR] ", prefix, err)
		return
	}

	number := rsc.DeletePodWithNumber(replica, int(replica.Status.Replicas), prefix)

	if number != int(replica.Status.Replicas) {
		fmt.Println("[ERROR] ", prefix, "only delete ", number, " pods. remain ",
			int(replica.Status.Replicas)-number, " pods.")
	} else {
		fmt.Println(prefix, "Successfully delete ", number, " pods in replicaset: ", replica.Name)
	}
}

func (rsc *ReplicaSetController) AddPod(event tool.Event) {
	prefix := "[ReplicaSet] [AddPod] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	fmt.Println(prefix, "event.Val: ", event.Val)
	rsc.PodInformer.Set(event.Key, event.Val)

	// add pod后首先等scheduler调度,因此不进行操作，操作在UpdatePod中进行
}

func (rsc *ReplicaSetController) UpdatePod(event tool.Event) {
	prefix := "[ReplicaSet] [UpdatePod] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	fmt.Println(prefix, "event.Val: ", event.Val)
	rsc.PodInformer.Set(event.Key, event.Val)

	pod := &core.Pod{}
	err := json.Unmarshal([]byte(event.Val), pod)
	if err != nil {
		fmt.Println("[ERROR] ", prefix, "pod Unmarshal", err)
		return
	}

	if pod.Status.Phase != "Running" {
		fmt.Println(prefix, "phase is not satisfied(only need Running):", pod.Status.Phase)
		return
	}

	replica_cache := rsc.ReplicasetInformer.GetCache()
	for _, value := range *replica_cache {
		replica := &core.ReplicaSet{}
		err := json.Unmarshal([]byte(value), replica)
		if err != nil {
			log.Println("[ERROR] ", prefix, "replicaset Unmarshal", err)
			return
		}
		if Match(replica, pod) {
			fmt.Println(prefix, "find replicaset match pod: ", replica.Name)
			pod.OwnerReferences = append(pod.OwnerReferences,
				v1.OwnerReference{
					Name: replica.Name,
					UID:  replica.UID,
				})
			replica.Status.Replicas++
			updateReplicaSet(rsc, replica, event.Key, prefix)
			return
		}
	}
	fmt.Println(prefix, "not find replicaset match pod: ", pod.Name)
}

func (rsc *ReplicaSetController) DeletePod(event tool.Event) {
	prefix := "[ReplicaSet] [DeletePod] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	fmt.Println(prefix, "event.Val: ", event.Val)
	rsc.PodInformer.Delete(event.Key)
}

func (rsc *ReplicaSetController) CreatePodWithNumber(replica *core.ReplicaSet, number int, prefix string) int {
	err_num := 0
	for i := 0; i < number; i++ {
		pod := &core.Pod{
			ObjectMeta: replica.Spec.Template.ObjectMeta,
			Spec:       replica.Spec.Template.Spec,
			Status: core.PodStatus{
				Phase: "Pending",
			},
		}
		pod.Name = replica.Name + "-" + random.GenerateRandomString(5)
		pod.OwnerReferences = append(pod.OwnerReferences,
			v1.OwnerReference{
				Name: replica.Name,
				UID:  replica.UID,
			})

		fmt.Println(prefix, "crate pod: ", pod.Name)
		err := create.CreatePod(pod)
		if err != nil {
			fmt.Println("[ERROR] ", prefix, "create pod error: ", err)
			err_num++
		}
		time.Sleep(time.Second * 4)
	}
	return number - err_num
}

func (rsc *ReplicaSetController) DeletePodWithNumber(replica *core.ReplicaSet, number int, prefix string) int {
	pod_cache := rsc.PodInformer.GetCache()
	success := 0

	for _, value := range *pod_cache {
		pod := &core.Pod{}
		err := json.Unmarshal([]byte(value), pod)
		if err != nil {
			log.Println("[ERROR] ", prefix, err)
			continue
		}
		for _, owner := range pod.OwnerReferences {
			if owner.UID == replica.UID {
				fmt.Println(prefix, "find pod: ", pod.Name)
				values := url.Values{}
				values.Add("PodName", pod.Name)
				err := web.SendHttpRequest("DELETE", apiconfig.Server_URL+apiconfig.POD_PATH+"?"+values.Encode(),
					web.WithPrefix(prefix),
					web.WithLog(true))
				if err != nil {
					fmt.Println("[ERROR] ", prefix, err)
					continue
				}
				success++
				if success == number {
					return success
				}
				time.Sleep(time.Second * 4)
			}
		}
	}
	return success
}

func (rsc *ReplicaSetController) worker() {
	prefix := "[ReplicaSet] [worker] "
	for {
		//TODO: optmize can use channel or condition variable
		if !rsc.queue.IsEmpty() {
			key := rsc.queue.Pop().(string)
			val, exist := (*rsc.ReplicasetInformer.GetCache())[key]
			if !exist {
				fmt.Println("[ERROR] ", prefix, "cache doesn't have key:", key)
				continue
			}
			replica := &core.ReplicaSet{}
			err := json.Unmarshal([]byte(val), replica)
			if err != nil {
				fmt.Println("[ERROR] ", prefix, err)
				return
			}

			diff := *replica.Spec.Replicas - replica.Status.Replicas
			var num int

			if diff > 0 {
				// create new
				fmt.Println(prefix, "start to create %d pod(s).", diff)
				num = rsc.CreatePodWithNumber(replica, int(diff), prefix)
			} else if diff < 0 {
				// delete old
				fmt.Println(prefix, "start to delete %d pod(s).", -diff)
				num = rsc.DeletePodWithNumber(replica, int(-diff), prefix)
			} else {
				// do nothing
				fmt.Println(prefix, "do nothing.")
			}

			replica.Status.Replicas += int32(num)

			// 更新replicaset
			data, err := json.Marshal(replica)
			if err != nil {
				fmt.Println("[kubectl] [create] [RunCreateReplicaSet] failed to marshal:", err)
			}

			// 创建 PUT 请求
			err = web.SendHttpRequest("POST", apiconfig.Server_URL+apiconfig.REPLICASET_PATH,
				web.WithPrefix("[kubectl] [create] [RunCreateReplicaSet] "),
				web.WithBody(bytes.NewBuffer(data)),
				web.WithLog(true))
			if err != nil {
				return
			}

		} else {
			time.Sleep(time.Second * 3)
		}
	}
}

func (rsc *ReplicaSetController) Run() {
	go rsc.PodInformer.Run()
	go rsc.ReplicasetInformer.Run()
	go rsc.worker()
	select {}
}