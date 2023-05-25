package kube_service

import (
	"encoding/json"
	"fmt"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	kube_proxy "minik8s/configs"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
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
	res.DNSInformer.AddEventHandler(tool.Added, res.DNSUpdateHandler)
	res.DNSInformer.AddEventHandler(tool.Modified, res.DNSUpdateHandler)

	return res
}

func (sm *ServiceManager) DNSUpdateHandler(event tool.Event) {
	prefix := "[ServiceManager][DNSUpdateHandler]"
	//fmt.Println(prefix + "key:" + event.Key)
	dns := &core.DNS{}
	err := json.Unmarshal([]byte(event.Val), dns)
	if err != nil {
		fmt.Println(prefix + err.Error())
		return
	}
	if dns.Status == core.FileCreatedStatus {
		fmt.Println(prefix + "is file created status")
		err = tool.AddPod(GetGatewayPodSingleton(dns.Metadata.Name)) // create specific GatewayPod
		if err != nil {
			fmt.Println(prefix + err.Error())
			return
		}
		err = tool.UpdateService(GetGatewayServiceSingleton(dns.Metadata.Name)) // create Specific GatewayService
		if err != nil {
			fmt.Println(prefix + err.Error())
			return
		}
		sm.lock.Lock()
		sm.name2DNS[dns.Metadata.Name] = dns
		sm.lock.Unlock()
	}
	if dns.Status == core.DeletedStatus {
		fmt.Println(prefix + "is deleted status")
		err = tool.DeleteService(GetGatewayServiceSingleton(dns.Metadata.Name))
		if err != nil {
			fmt.Println(prefix + err.Error())
			return
		}
		err = tool.DeletePod(dns.Metadata.Name)
		if err != nil {
			fmt.Println(prefix + err.Error())
			return
		}
	}
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
		if dns != nil {
			if dns.ServiceMeta.Name == kube_proxy.CoreDNSServiceName { // service exists before
				fmt.Println(prefix + "CoreDNS exists!")
				return
			}
			fmt.Println(prefix + "DNS parse error!!")
		}
		// insert to etcd,  it will be handled by kubelet
		err = tool.AddPod(GetCoreDNSPodSingleton())
		if err != nil {
			fmt.Println(prefix + err.Error())
			return
		}
		time.Sleep(10 * time.Second)
		// create coreDNS service
		err = tool.UpdateService(GetCoreDNSServiceSingleton())
		time.Sleep(15 * time.Second)
		if err != nil {
			fmt.Println(prefix + err.Error())
			return
		}
	}
}

// should be run by go routine
func (sm *ServiceManager) checkNginx() {
	prefix := "[ServiceManager][DNSChecker]"
	for {
		time.Sleep(2 * time.Second)
		var rms []string
		sm.lock.Lock()
		if len(sm.name2DNS) > 0 {
			fmt.Println(prefix + "running")
		} else {
			sm.lock.Unlock()
			continue
		}
		for k, v := range sm.name2DNS {
			res, err := tool.GetService(kube_proxy.GatewayServicePrefix + k)
			if err != nil {
				fmt.Println(prefix + err.Error())
				continue
			}
			if res == nil {
				continue
			}
			if res.ServiceSpec.Status.Phase == service.ServiceRunningPhase { // the nginx is running
				fmt.Println(prefix + "key=" + k + "nginx service created!")
				v.Status = core.ServiceCreatedStatus
				v.Spec.GatewayIp = res.ServiceSpec.ClusterIP
				err = tool.UpdateDNS(v)
				if err != nil {
					fmt.Println(prefix + err.Error())
					continue
				}
				rms = append(rms, k)
			}
		}
		for _, v := range rms {
			delete(sm.name2DNS, v)
		}
		sm.lock.Unlock()
	}
}
func (sm *ServiceManager) Run() {
	// 启动ServiceInformer, it can listen and create/run service
	go sm.ServiceInformer.Run()
	// 启动 dns 注册服务
	go InitCoreDNS()
	// 启动 serverManager端的 DNS informer，主要和service配合，修改gatewayIp
	go sm.DNSInformer.Run()
	// 每2s 检查 nginx的状态
	go sm.checkNginx()
}
