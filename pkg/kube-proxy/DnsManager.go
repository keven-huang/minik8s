package kube_proxy

import (
	"encoding/json"
	"fmt"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
)

func NewDnsManager() *DNSManager {
	res := &DNSManager{}
	res.Key2Dns = make(map[string]*core.DNS)
	res.isDead = false
	res.DNSInformer = informer.NewInformer(apiconfig.DNS_PATH)
	res.Register() // register handler
	return res
}

func (DNSManager *DNSManager) UpdateDNSHandler(event tool.Event) {
	prefix := "[DNSManager][UpdateDns]"
	fmt.Println(prefix + "key:" + event.Key)
	dm := &core.DNS{}
	err := json.Unmarshal([]byte(event.Val), dm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, ok := DNSManager.Key2Dns[event.Key]
	if !ok { // first
		// TODO
	} else { // update
		// TODO
	}
}

func (DNSManager *DNSManager) GetDNSHandler(event tool.Event) {

}

func (DNSManager *DNSManager) DeleteDNSHandler(event tool.Event) {

}

func (DNSManager *DNSManager) Register() {
	DNSManager.DNSInformer.AddEventHandler(tool.Added, DNSManager.UpdateDNSHandler)
	DNSManager.DNSInformer.AddEventHandler(tool.Modified, DNSManager.UpdateDNSHandler)
	DNSManager.DNSInformer.AddEventHandler(tool.Deleted, DNSManager.DeleteDNSHandler)
}
