package kube_service

import (
	"fmt"
	kube_proxy "minik8s/configs"
	"minik8s/pkg/api/core"
	"minik8s/pkg/service"
	myJson "minik8s/pkg/util/json"
)

var coreDnsPod *core.Pod
var coreDnsService *service.Service
var gatewayPod *core.Pod            // template
var gatewayService *service.Service // template

func GetGatewayPodSingleton(name string) *core.Pod {
	prefix := "[Singletons][GatewayPod]"
	resPod := &core.Pod{}
	fmt.Println(prefix + "in")
	if gatewayPod == nil {
		gatewayPod = &core.Pod{}
		err := myJson.GetFromYaml(kube_proxy.GatewayPodYamlPath, gatewayPod)
		if err != nil {
			fmt.Println(prefix + err.Error())
			return nil
		}
	}
	tmpSpec := gatewayPod.Spec // deep copy
	resPod.Spec = tmpSpec
	resPod.Spec.Volumes[0].HostPath = kube_proxy.NginxPrefix + "/" + name
	resPod.Spec.Containers[0].Name = kube_proxy.GatewayContainerPrefix + name
	resPod.ObjectMeta.Name = kube_proxy.GatewayPodPrefix + name
	//resPod.Name = kube_proxy.GatewayPodPrefix + name
	tmpLabels := gatewayPod.Labels
	resPod.Labels = tmpLabels
	resPod.Labels["dnsName"] = name // pod label for select
	return resPod
}

func GetGatewayServiceSingleton(dns *core.DNS) *service.Service {
	prefix := "[Singletons][GatewayService]"
	fmt.Println(prefix + "in")
	resService := &service.Service{}
	if gatewayService == nil { // create
		gatewayService = &service.Service{}
		err := myJson.GetFromYaml(kube_proxy.GatewayServiceYamlPath, gatewayService)
		if err != nil {
			fmt.Println(prefix + err.Error())
			return nil
		}
	}
	tmpMeta := gatewayService.ServiceMeta
	tmpMeta.Name = kube_proxy.GatewayServicePrefix + dns.Metadata.Name
	tmpSpec := gatewayService.ServiceSpec
	tmpSpec.Selector["dnsName"] = dns.Metadata.Name // this should be matched with pod's label
	tmpSpec.ClusterIP = dns.Spec.GatewayIp
	resService.ServiceSpec = tmpSpec
	resService.ServiceMeta = tmpMeta
	return resService
}

func GetCoreDNSPodSingleton() *core.Pod {
	prefix := "[Singletons][CoreDNSPod]"
	if coreDnsPod == nil { // create
		coreDnsPod = &core.Pod{}
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
		coreDnsService = &service.Service{}
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
