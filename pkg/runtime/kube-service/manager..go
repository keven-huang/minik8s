package kube_service

import (
	"encoding/json"
	"fmt"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	"minik8s/pkg/service"
)

type ServiceManager struct {
	// ServiceName->RuntimeService
	ServiceMapping map[string]*RuntimeService
	// Informer
	ServiceInformer informer.Informer
}

func NewServiceManager() *ServiceManager {
	res := &ServiceManager{}
	res.ServiceMapping = make(map[string]*RuntimeService)
	// 初始化serviceInformer
	res.ServiceInformer = informer.NewInformer(apiconfig.SERVICE_PATH)
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
		fmt.Println("[kube-service][manager][deleteServiceHandler]" + event.Key)
		var name string
		err := json.Unmarshal([]byte(event.Val), &name)
		if err != nil {
			return
		}
		lastService, ok := res.ServiceMapping[name]
		if !ok {
			fmt.Println("[warn]:" + "fail to delete service " + name)
			return
		} else {
			lastService.Delete()             // delete service
			delete(res.ServiceMapping, name) // delete from map
		}
	})

	return res
}

func (sm *ServiceManager) Run() {
	// 启动ServiceInformer, it can add service
	go sm.ServiceInformer.Run()
}
