package kube_service

import (
	"fmt"
	"minik8s/pkg/api/core"
	kube_proxy "minik8s/pkg/kube-proxy"
	"minik8s/pkg/service"
	myJson "minik8s/pkg/util/json"
)

var coreDnsPod *core.Pod
var coreDnsService *service.Service

func GetCoreDNSPodSingleton() *core.Pod {
	prefix := "[Singletons][CoreDNSPod]"
	if coreDnsPod == nil { // create
		err := myJson.GetFromYaml(kube_proxy.CoreDnsPodYamlPath, coreDnsPod)
		if err != nil {
			fmt.Println(prefix + err.Error())
			return nil
		}
		return coreDnsPod
	} else {
		return coreDnsPod
	}
}

func GetCoreDNSServiceSingleton() *service.Service {
	prefix := "[Singletons][CoreDNSService]"
	if coreDnsService == nil { // create
		err := myJson.GetFromYaml(kube_proxy.CoreDnsServiceYamlPath, coreDnsService)
		if err != nil {
			fmt.Println(prefix + err.Error())
			return nil
		}
		return coreDnsService
	} else {
		return coreDnsService
	}
}
