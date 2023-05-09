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
	res.ServiceInformer = informer.NewInformer(apiconfig.SERVICE_PATH)
	// 设置watch add service的回调函数
	res.ServiceInformer.AddEventHandler(tool.Added, func(event tool.Event) {
		fmt.Println("[info]:" + "addHandler" + event.Key)
		newService := &service.Service{}
		err := json.Unmarshal([]byte(event.Val), newService)
		if err != nil {
			return
		}
		lastService, ok := res.ServiceMapping[newService.ServiceMeta.Name]
		if !ok { // 新建
			res.ServiceMapping[newService.ServiceMeta.Name] = CreateService(newService, &res.ServiceInformer)
		} else { // 删除老的，创建新的
			lastService.Delete()
			delete(res.ServiceMapping, newService.ServiceMeta.Name)
			res.ServiceMapping[newService.ServiceMeta.Name] = CreateService(newService, &res.ServiceInformer)
		}
	})
	return res
}

func (sm *ServiceManager) Run() {
	// 启动ServiceInformer
	go sm.ServiceInformer.Run()

}
