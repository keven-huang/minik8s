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
	"minik8s/pkg/kubelet/monitor"
	"minik8s/pkg/util/file"
	"minik8s/pkg/util/web"
	"regexp"
	"time"
)

type Kubelet struct {
	Monitor     *monitor.Monitor
	PodInformer informer.Informer
	node        core.Node
}

func NewKubelet(name string, nodeIp string, masterIp string) (*Kubelet, error) {
	node := core.Node{}
	node.Name = name
	node.Spec.NodeIP = nodeIp
	apiconfig.Server_URL = masterIp
	err := tool.AddNode(&node)
	if err != nil {
		return nil, err
	}
	return &Kubelet{
		Monitor:     monitor.NewMonitor(9400),
		PodInformer: informer.NewInformer(apiconfig.POD_PATH),
		node:        node,
	}, nil
}

func (k *Kubelet) Register() {
	k.PodInformer.AddEventHandler(tool.Added, k.CreatePod)
	k.PodInformer.AddEventHandler(tool.Modified, k.UpdatePod)
	k.PodInformer.AddEventHandler(tool.Deleted, k.DeletePod)
}

func (k *Kubelet) CreatePod(event tool.Event) {
	prefix := "[kubelet] [CreatePod] "
	// handle event
	fmt.Println(prefix, "event.Type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	fmt.Println(prefix, "event.Val: ", event.Val)
	k.PodInformer.Set(event.Key, event.Val)

}

func (k *Kubelet) UpdatePod(event tool.Event) {
	prefix := "[kubelet] [UpdatePod] "
	// handle event
	fmt.Println(prefix, "event.Type: ", tool.GetTypeName(event))
	fmt.Println(prefix, "event.Key: ", event.Key)
	fmt.Println(prefix, "event.Val: ", event.Val)
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

	// 判断创建pod是否是gpu类型
	if pod.Spec.GPUJob == true {
		// 获取gpu运行文件至job目录
		fmt.Println(prefix, "get gpufiles")
		err = GetGpuJobFile(pod.Spec.GPUJobName)
		if err != nil {
			fmt.Println(prefix, "[ERROR]", err)
			return
		}
	}

	metaData, netSetting, err := dockerClient.CreatePod(*pod)
	if err != nil {
		fmt.Println(prefix, err)
		return
	}

	// 创建成功 修改Status
	pod.Status.Phase = core.PodRunning

	// 更新podIp
	pod.Status.PodIP = netSetting.IPAddress

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
	key := k.PodInformer.Get(event.Key)
	fmt.Println("event.Val: ", key)
	k.PodInformer.Delete(event.Key)

	pod := &core.Pod{}
	err := json.Unmarshal([]byte(key), pod)
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

func GetGpuJobFile(jobname string) error {
	// 获取gpu运行文件至job目录
	jobFile := tool.GetJobFile(jobname)
	if len(jobFile.Program) == 0 || len(jobFile.Slurm) == 0 {
		return fmt.Errorf("[kubelet] [UpdatePod] get gpufiles empty")
	}
	// get program
	program_name := jobFile.JobName + ".cu"
	err := file.MakeFile(jobFile.Program, program_name, apiconfig.JOB_FILE_DIR_PATH)
	if err != nil {
		return err
	}
	// get slurm
	slurm_name := jobFile.JobName + ".slurm"
	err = file.MakeFile(jobFile.Slurm, slurm_name, apiconfig.JOB_FILE_DIR_PATH)
	if err != nil {
		return err
	}
	return nil
}

func (k *Kubelet) Listener(needLog bool) {
	re := regexp.MustCompile(`^Exited \((\d+)\)`)

	// docker ps every 15 seconds
	for {
		containers, err := dockerClient.GetAllContainers()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, con := range containers {
			match := re.FindStringSubmatch(con.Status)
			Exited := false
			ExitedValue := ""
			if len(match) > 1 {
				if needLog {
					fmt.Printf("%v Exited:%s\n", con.Names[0][1:], match[1])
				}
				Exited = true
				ExitedValue = match[1]
			} else {
				if needLog {
					fmt.Printf("%v not found.\n", con.Names[0][1:])
				}
			}
			if Exited {
				pod_cache := k.PodInformer.GetCache()
				for _, v := range *pod_cache {
					pod := &core.Pod{}
					err := json.Unmarshal([]byte(v), pod)
					if err != nil {
						fmt.Println(err)
						return
					}
					for _, c := range pod.Spec.Containers {
						if c.Name == con.Names[0][1:] && pod.Status.Phase == "Running" {
							fmt.Println("pod:", pod.Name, " container:", c.Name, " Exited With", ExitedValue)
							if ExitedValue == "0" {
								pod.Status.Phase = "Succeeded"
							} else {
								pod.Status.Phase = "Failed"
							}

							data, err := json.Marshal(pod)
							if err != nil {
								fmt.Println(err)
								return
							}
							err = web.SendHttpRequest("POST", apiconfig.Server_URL+apiconfig.POD_PATH,
								web.WithPrefix("[kubelet] [Listener]"), web.WithBody(bytes.NewBuffer(data)))
							if err != nil {
								fmt.Println(err)
								return
							}
							if needLog {
								fmt.Printf("%v update status to %s.\n", con.Names[0][1:], pod.Status.Phase)
							}
							break
						}
					}

				}
			}
		}
		if needLog {
			fmt.Println("-----------")
		}
		time.Sleep(15 * time.Second)
	}
}

func (k *Kubelet) Run() {
	go k.PodInformer.Run()
	go k.Listener(true)
	select {}
}
