package kube_service

import (
	"encoding/json"
	"fmt"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	kube_proxy "minik8s/pkg/kube-proxy"
	"minik8s/pkg/service"
	"strings"
	"sync"
	"time"
)

type ServiceManager struct {
	// ServiceName->RuntimeService
	// used for delete
	ServiceMapping map[string]*RuntimeService
	// Informer for service
	ServiceInformer informer.Informer
	// Informer for dns
	DNSInformer informer.Informer
	// mapping from dns-config-name -> dns data structure
	// used for delete...
	name2DNS map[string]*core.DNS
	// lock
	lock sync.RWMutex
}

func NewServiceManager() *ServiceManager {
	res := &ServiceManager{}
	var lock sync.RWMutex
	res.lock = lock
	res.ServiceMapping = make(map[string]*RuntimeService)
	res.name2DNS = make(map[string]*core.DNS)
	// 初始化serviceInformer
	res.ServiceInformer = informer.NewInformer(apiconfig.SERVICE_PATH)
	// 初始化 DNSInformer
	res.DNSInformer = informer.NewInformer(apiconfig.DNS_PATH)

	// 设置watch add service event 的回调函数
	res.ServiceInformer.AddEventHandler(tool.Added, func(event tool.Event) {
		fmt.Println("[kube-service][manage][addServiceHandler]:" + event.Key)
		newService := &service.Service{}
		err := json.Unmarshal([]byte(event.Val), newService)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		lastService, ok := res.ServiceMapping[newService.ServiceMeta.Name]
		if !ok { // 新建
			res.ServiceMapping[newService.ServiceMeta.Name] = CreateService(newService)
		} else { // 删除老的，创建新的
			lastService.Delete()                                    // delete service
			delete(res.ServiceMapping, newService.ServiceMeta.Name) // delete from map
			res.ServiceMapping[newService.ServiceMeta.Name] = CreateService(newService)
		}
	})
	// 设置watch delete event的回调
	res.ServiceInformer.AddEventHandler(tool.Deleted, func(event tool.Event) {
		// delete by name
		prefix := "[kube-service][manager][deleteServiceHandler]"
		fmt.Println(prefix + event.Key)
		strs := strings.Split(event.Key, "/")
		var name string
		name = strs[4]
		fmt.Println(prefix + name)
		lastService, ok := res.ServiceMapping[name]
		if !ok {
			fmt.Println(prefix + "fail to find service " + name)
			return
		} else {
			lastService.Delete()             // delete service
			delete(res.ServiceMapping, name) // delete from map
		}
	})

	// TODO DNS-informer
	return res
}

// should be run by go-routine
func InitCoreDNS() {
	for {
		prefix := "[serviceManager][initCoreDNS]"
		dns, err := tool.GetService(kube_proxy.CoreDNSServiceName)
		if err != nil {
			fmt.Println(prefix + err.Error())
			return
		}
		if dns != nil { // service exists before
			fmt.Println(prefix + "CoreDNS exists!")
			return
		}
		// insert to etcd,  it will be handled by kubelet
		err = tool.UpdatePod(GetCoreDNSPodSingleton())
		if err != nil {
			fmt.Println(prefix + err.Error())
			return
		}
		time.Sleep(5 * time.Second)
		// create coreDNS service
		err = tool.UpdateService(GetCoreDNSServiceSingleton())
		if err != nil {
			fmt.Println(prefix + err.Error())
			return
		}
	}
}

func (sm *ServiceManager) Run() {
	// 启动ServiceInformer, it can listen and create/run service
	go sm.ServiceInformer.Run()
	// 启动 dns 注册服务
	go InitCoreDNS()
}
