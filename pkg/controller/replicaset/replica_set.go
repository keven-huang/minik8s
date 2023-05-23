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

func MakeOwnerReference(replica *core.ReplicaSet) v1.OwnerReference {
	return v1.OwnerReference{
		APIVersion: "apps/v1",
		Kind:       "ReplicaSet",
		Name:       replica.Name,
		UID:        replica.UID,
	}
}

func NewReplicaSetController() (*ReplicaSetController, error) {
	rsc := &ReplicaSetController{
		PodInformer:        informer.NewInformer(apiconfig.POD_PATH),
		ReplicasetInformer: informer.NewInformer(apiconfig.REPLICASET_PATH),
		queue:              q.NewConcurrentQueue(),
	}

	prefix := "[ReplicaSetController] [NewReplicaSetController] "
	replica_cache := rsc.ReplicasetInformer.GetCache()
	pod_cache := rsc.PodInformer.GetCache()

	for key, value := range *replica_cache {
		replica := &core.ReplicaSet{}
		err := json.Unmarshal([]byte(value), replica)
		if err != nil {
			log.Println("[ERROR] ", "replicaset Unmarshal", err)
			continue
		}
		fmt.Println(replica.Name)

		var Replicas = 0
		for pod_key, pod := range *pod_cache {
			p := &core.Pod{}
			err := json.Unmarshal([]byte(pod), p)
			if err != nil {
				log.Println("[ERROR] ", "pod Unmarshal", err)
				continue
			}
			if Match(replica, p) {
				Replicas++
				flag := false
				for _, own := range p.OwnerReferences {
					if own.Name == replica.Name {
						flag = true
						break
					}
				}
				if !flag {
					p.OwnerReferences = append(p.OwnerReferences, MakeOwnerReference(replica))
					err = updatePodToServer(rsc, p, pod_key, prefix)
					if err != nil {
						log.Println("[ERROR] ", prefix, "pod update owner", err)
						continue
					}
				}
			}
		}

		if int32(Replicas) != replica.Status.Replicas {
			replica.Status.Replicas = int32(Replicas)
			updateReplicaSetToServer(rsc, replica, key, prefix)
			rsc.queue.Push(key)
		}
	}

	return rsc, nil
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
		v, exists := pod.ObjectMeta.Labels[key]
		if !exists || v != value {
			return false
		}
	}
	for _, expression := range rs.Spec.Selector.MatchExpressions {
		// Valid operators are In, NotIn, Exists and DoesNotExist.
		key := expression.Key
		val, exists := pod.ObjectMeta.Labels[key]
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

func updateReplicaSetToServer(rsc *ReplicaSetController, replica *core.ReplicaSet, key string, prefix string) error {
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

	rsc.ReplicasetInformer.Set(key, string(data))
	return nil
}

func updatePodToServer(rsc *ReplicaSetController, pod *core.Pod, key string, prefix string) error {
	data, err := json.Marshal(pod)
	if err != nil {
		fmt.Println(prefix, "failed to marshal:", err)
		return err
	}

	err = web.SendHttpRequest("POST", apiconfig.Server_URL+apiconfig.POD_PATH,
		web.WithPrefix(prefix), web.WithBody(bytes.NewBuffer(data)))
	if err != nil {
		return err
	}

	rsc.PodInformer.Set(key, string(data))
	return nil
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
			pod.OwnerReferences = append(pod.OwnerReferences, MakeOwnerReference(replica))
			replica.Status.Replicas++
			flag = true
			fmt.Println(prefix, "find pod: ", pod.Name)
		}
	}

	if flag {
		updateReplicaSetToServer(rsc, replica, event.Key, prefix)
	}

	if *replica.Spec.Replicas != replica.Status.Replicas {
		rsc.queue.Push(event.Key)
	}

}

// 主要是更新Status.Replicas,不进行操作，由worker进行操作
func (rsc *ReplicaSetController) UpdateReplicaset(event tool.Event) {
	prefix := "[ReplicaSet] [UpdateReplicaset] "
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

	if replica.Status.Replicas != *replica.Spec.Replicas {
		rsc.queue.Push(event.Key)
	}
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

// add pod后首先等scheduler调度,因此不进行操作，操作在UpdatePod中进行
func (rsc *ReplicaSetController) AddPod(event tool.Event) {
	prefix := "[ReplicaSet] [AddPod] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	fmt.Println(prefix, "event.Val: ", event.Val)
	rsc.PodInformer.Set(event.Key, event.Val)

}

// 可能调用的情况：
// 1. scheduler调度成功，更新pod中的NodeName
// 2. kebelet创建成功，pod的phase变为Running，此时需要判断是否匹配replicaset并加入
// 3. kubelet发现container失败或者成功，更新pod的状态
// 综上,只关注phase为Running和owner中没有ReplicaSet的pod
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

	// 判断已经有owner的情况，也就是说本来就是replicaSet自己create的
	if len(pod.OwnerReferences) != 0 {
		for _, owner := range pod.OwnerReferences {
			if owner.Kind == "ReplicaSet" {
				fmt.Println(prefix, "already has ReplicaSet owner, ", owner.Name)
				return
			}
		}
	}

	// 默认只对应一个replicaset
	replica_cache := rsc.ReplicasetInformer.GetCache()
	for replicaset_key, value := range *replica_cache {
		replica := &core.ReplicaSet{}
		err := json.Unmarshal([]byte(value), replica)
		if err != nil {
			log.Println("[ERROR] ", prefix, "replicaset Unmarshal", err)
			return
		}
		if Match(replica, pod) {
			fmt.Println(prefix, "find replicaset match pod: ", replica.Name)
			pod.OwnerReferences = append(pod.OwnerReferences, MakeOwnerReference(replica))
			err := updatePodToServer(rsc, pod, event.Key, prefix)
			if err != nil {
				fmt.Println("[ERROR] ", prefix, "updatePodToServer failed ", err)
				return
			}

			replica.Status.Replicas++
			err = updateReplicaSetToServer(rsc, replica, replicaset_key, prefix)
			if err != nil {
				fmt.Println("[ERROR] ", prefix, "updateReplicaSetToServer failed ", err)
				return
			}
			rsc.queue.Push(replicaset_key)
			return
		}
	}
	fmt.Println(prefix, "not find replicaset match pod: ", pod.Name)
}

// 主要看被删除的Pod是否有ReplicaSet的owner，如果是则更新ReplicaSet
func (rsc *ReplicaSetController) DeletePod(event tool.Event) {
	prefix := "[ReplicaSet] [DeletePod] "
	fmt.Println(prefix, "event.type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	key := rsc.PodInformer.Get(event.Key)
	rsc.PodInformer.Delete(event.Key)

	pod := &core.Pod{}
	err := json.Unmarshal([]byte(key), pod)
	if err != nil {
		fmt.Println("[ERROR] ", prefix, "pod Unmarshal failed ", err)
		return
	}

	// 判断是否有owner
	if len(pod.OwnerReferences) == 0 {
		fmt.Println(prefix, "pod has no owner.")
		return
	}

	// 判断是否有ReplicaSet的owner
	flag := false
	Name := ""
	for _, owner := range pod.OwnerReferences {
		if owner.Kind == "ReplicaSet" {
			fmt.Println(prefix, "pod has ReplicaSet owner, ", owner.Name)
			flag = true
			Name = owner.Name
			break
		}
	}
	if !flag {
		fmt.Println(prefix, "pod has no ReplicaSet owner.")
		return
	}

	for replicasetKey, val := range *rsc.ReplicasetInformer.GetCache() {
		replica := &core.ReplicaSet{}
		err := json.Unmarshal([]byte(val), replica)
		if err != nil {
			fmt.Println("[ERROR] ", prefix, "replica Unmarshal failed ", err)
			return
		}
		if replica.Name == Name {
			fmt.Println(prefix, "find replica match pod: ", replica.Name)
			replica.Status.Replicas--
			err := updateReplicaSetToServer(rsc, replica, replicasetKey, prefix)
			if err != nil {
				fmt.Println("[ERROR] ", prefix, "update replica failed ", err)
				return
			}
			rsc.queue.Push(replicasetKey)
			return
		}
	}

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
		pod.OwnerReferences = append(pod.OwnerReferences, MakeOwnerReference(replica))

		fmt.Println(prefix, "crate pod: ", pod.Name)
		err := create.CreatePod(pod)
		if err != nil {
			fmt.Println("[ERROR] ", prefix, "create pod error: ", err)
			err_num++
		}
		//time.Sleep(time.Second * 3)
	}
	return number - err_num
}

// DeletePodWithNumber return the number of pods deleted
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
				fmt.Println(prefix, "delete pod: ", pod.Name)
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
				//time.Sleep(time.Second * 3)
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

			for *replica.Spec.Replicas-replica.Status.Replicas != 0 {
				diff := *replica.Spec.Replicas - replica.Status.Replicas
				var num int

				if diff > 0 {
					// create new
					fmt.Print(prefix, "start to create %d pod(s).\n", diff)
					num = rsc.CreatePodWithNumber(replica, 1, prefix)
				} else if diff < 0 {
					// delete old
					fmt.Print(prefix, "start to delete %d pod(s).\n", -diff)
					num = -rsc.DeletePodWithNumber(replica, 1, prefix)
				} else {
					// do nothing
					fmt.Println(prefix, "do nothing.")
				}

				fmt.Print(prefix, "finish to create/delete %d pod(s).\n", num)

				if num != 0 {
					replica.Status.Replicas += int32(num)
				}
			}

			// 更新replicaset
			data, err := json.Marshal(replica)
			if err != nil {
				fmt.Println(prefix, "failed to marshal:", err)
			}

			// 创建 PUT 请求
			err = web.SendHttpRequest("POST", apiconfig.Server_URL+apiconfig.REPLICASET_PATH,
				web.WithPrefix(prefix),
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
