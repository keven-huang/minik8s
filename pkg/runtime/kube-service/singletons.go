package kube_service

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	kube_proxy "minik8s/configs"
	"minik8s/pkg/api/core"
	"minik8s/pkg/service"
	"os"
)

var coreDnsPod *core.Pod
var coreDnsService *service.Service
var gatewayPod *core.Pod            // template
var gatewayService *service.Service // template

func GetFromYaml(filename string, a interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open file %s", filename)
	}
	defer file.Close()
	// Read the YAML file
	var dataMap map[string]interface{}

	yamlFile, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read yaml error")
	}

	err = yaml.Unmarshal(yamlFile, &dataMap)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// get kind
	var kind string
	kind = dataMap["kind"].(string)
	fmt.Println("kind:", kind)
	switch kind {
	case "Pod":
		//pod := &core.Pod{}
		err := yaml.Unmarshal(yamlFile, a.(*core.Pod))
		data, err := json.Marshal(a.(*core.Pod))
		fmt.Println(string(data))
		if err != nil {
			return err
		}
		return nil
	case "Service":
		//s := &service.Service{}
		err := yaml.Unmarshal(yamlFile, a.(*service.Service))
		if err != nil {
			return err
		}
		return nil
	case "dns":
		//r := &core.DNS{}
		err := yaml.Unmarshal(yamlFile, a.(*core.DNS))
		if err != nil {
			return err
		}
		return nil
	}

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetGatewayPodSingleton(name string) *core.Pod {
	prefix := "[Singletons][GatewayPod]"
	resPod := &core.Pod{}
	fmt.Println(prefix + "in")
	if gatewayPod == nil {
		gatewayPod = &core.Pod{}
		err := GetFromYaml(kube_proxy.GatewayPodYamlPath, gatewayPod)
		if err != nil {
			fmt.Println(prefix + err.Error())
			return nil
		}
	}
	tmpSpec := gatewayPod.Spec // deep copy
	resPod.Spec = tmpSpec
	// TODO support many volumes
	resPod.Spec.Volumes[0].Name = kube_proxy.GatewayVolumePrefix + name
	resPod.Spec.Volumes[0].HostPath = kube_proxy.NginxPrefix + "/" + name
	resPod.Spec.Containers[0].VolumeMounts[0].Name = kube_proxy.GatewayVolumePrefix + name
	resPod.Spec.Containers[0].Name = kube_proxy.GatewayContainerPrefix + name
	resPod.ObjectMeta.Name = kube_proxy.GatewayPodPrefix + name
	//resPod.Name = kube_proxy.GatewayPodPrefix + name
	tmpLabels := gatewayPod.Labels
	resPod.Labels = tmpLabels
	resPod.Labels["dnsName"] = name // pod label for select
	resPod.Spec.NodeName = "node1"
	return resPod
}

func GetGatewayServiceSingleton(dns *core.DNS) *service.Service {
	prefix := "[Singletons][GatewayService]"
	fmt.Println(prefix + "in")
	resService := &service.Service{}
	if gatewayService == nil { // create
		gatewayService = &service.Service{}
		err := GetFromYaml(kube_proxy.GatewayServiceYamlPath, gatewayService)
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
		err := GetFromYaml(kube_proxy.CoreDnsPodYamlPath, coreDnsPod)
		data, err := json.Marshal(coreDnsPod)
		fmt.Println(string(data))
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
		err := GetFromYaml(kube_proxy.CoreDnsServiceYamlPath, coreDnsService)
		data, err := json.Marshal(coreDnsPod)
		fmt.Println(string(data))
		if err != nil {
			fmt.Println(prefix + err.Error())
			return nil
		}
		return coreDnsService
	} else {
		return coreDnsService
	}
}
