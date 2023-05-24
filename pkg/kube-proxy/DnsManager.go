package kube_proxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	"minik8s/pkg/client/tool"
	"minik8s/pkg/kubelet/dockerClient"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func NewDnsManager() *DNSManager {
	res := &DNSManager{}
	res.Key2Dns = make(map[string]*core.DNS)
	res.isDead = false
	res.DNSInformer = informer.NewInformer(apiconfig.DNS_PATH) // register the dns resource
	res.Register()                                             // register handler
	return res
}

func (DNSManager *DNSManager) UpdateDNSHandler(event tool.Event) {
	prefix := "[DNSManager][UpdateDns]"
	fmt.Println(prefix + "key:" + event.Key)
	dns := &core.DNS{}
	err := json.Unmarshal([]byte(event.Val), dns)
	if err != nil {
		fmt.Println(prefix + err.Error())
		return
	}
	switch dns.Status {
	case "": // first create
		{
			DNSManager.Key2Dns[event.Key] = dns
			DNSManager.copyDir(dns.Metadata.Name)
			// write nginx
			DNSManager.writeNginx(dns)
			dns.Status = core.FileCreatedStatus
			err := tool.UpdateDNS(dns) // updateETCD
			if err != nil {
				fmt.Println(prefix + err.Error())
				return
			}
			break
		}
	case core.ServiceCreatedStatus:
		{
			DNSManager.Key2Dns[event.Key] = dns
			DNSManager.writeNginx(dns)
			DNSManager.reloadNginx(dns)
			DNSManager.writeCoreDNS()
			break
		}
	}
}

func (DNSManager *DNSManager) DeleteDNSHandler(event tool.Event) {
	prefix := "[DNSManager][DeleteDNSHandler]"
	fmt.Println(prefix + "key: " + event.Key)
	dns, ok := DNSManager.Key2Dns[event.Key]
	if !ok {
		return
	} else {
		delete(DNSManager.Key2Dns, event.Key)
		DNSManager.writeCoreDNS()
		DNSManager.deleteDir(dns.Metadata.Name)
	}
}

func (DNSManager *DNSManager) Register() {
	DNSManager.DNSInformer.AddEventHandler(tool.Added, DNSManager.UpdateDNSHandler)
	DNSManager.DNSInformer.AddEventHandler(tool.Modified, DNSManager.UpdateDNSHandler)
	DNSManager.DNSInformer.AddEventHandler(tool.Deleted, DNSManager.DeleteDNSHandler)
}

func (DNSManager *DNSManager) writeNginx(dns *core.DNS) {
	prefix := "[DNSManager][WriteNginx]"
	var data []string
	data = append(data, "events {worker_connections 1024; }")
	data = append(data, "http {")
	data = append(data, "    server {", "         listen 80;")
	data = append(data, fmt.Sprintf("        server_name %s;", dns.Spec.Host)) // it can be deleted maybe
	data = append(data, DNSManager.generateConfig(dns)...)
	data = append(data, "        }")
	data = append(data, "}")
	file, err := os.OpenFile(NginxPrefix+"/"+dns.Metadata.Name+"/"+"nginx.conf", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(prefix + "open nginx file error" + err.Error())
		return
	}
	w := bufio.NewWriter(file)
	for _, v := range data {
		_, err := fmt.Fprintln(w, v)
		if err != nil {
			fmt.Println(prefix + "Write nginx data error" + err.Error())
		}
	}
	err = w.Flush()
	if err != nil {
		fmt.Println(prefix + err.Error())
	}
	return
}

// 生成nginx配置
func (DNSManager *DNSManager) generateConfig(dns *core.DNS) []string {
	var res []string
	for _, v := range dns.Spec.Paths {
		res = append(res, fmt.Sprintf("        location %s {", v.Name))
		res = append(res, fmt.Sprintf("        proxy_pass http://%s:%s;", v.Ip, v.Port))
		res = append(res, "        }")
	}
	return res
}

func (DNSManager *DNSManager) copyDir(Dns string) {
	prefix := "[DNSManager][CopyDir]"
	args := fmt.Sprintf("-r %s %s", NginxPath, NginxPrefix+"/"+Dns)
	res, err := execCommand("cp", args)
	if err != nil {
		fmt.Println(prefix + err.Error())
	} else {
		fmt.Println(res)
	}
}

func (DNSManager *DNSManager) deleteDir(Dns string) {
	prefix := "[DNSManager][deleteDir]"
	args := fmt.Sprintf("-rf %s", NginxPrefix+"/"+Dns)
	res, err := execCommand("rm", args)
	if err != nil {
		fmt.Println(prefix + err.Error())
	} else {
		fmt.Println(res)
	}
}

// add an entry to dns
// mapping from domainName -> gateway
// the gateway is a nginx
// TODO, file should be truncated?
func (DNSManager *DNSManager) writeCoreDNS() {
	prefix := "[DNSManager][WriteCoreDNS]"
	file, err := os.OpenFile(CoreDnsPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println(prefix + err.Error())
		return
	}
	writer := bufio.NewWriter(file)
	for _, v := range DNSManager.Key2Dns {
		if v.Status == core.ServiceCreatedStatus {
			cur := v.Spec.GatewayIp + " " + v.Spec.Host
			n, err := fmt.Fprintln(writer, cur)
			if err != nil {
				fmt.Println(prefix + err.Error())
			} else {
				fmt.Println(prefix + " success add " + strconv.Itoa(n) + "mappings")
			}
		}
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println(prefix + err.Error())
	}
	return
}

func (DNSManager *DNSManager) reloadNginx(dns *core.DNS) {
	prefix := "[DNSManager][reloadNginx]"
	cons, err := dockerClient.GetRunningContainers()
	if err != nil {
		fmt.Println(prefix + err.Error())
		return
	}
	var ids []string
	for _, val := range cons {
		// TODO names
		if strings.Contains(val.Names[0], GatewayContainerPrefix+dns.Metadata.Name) {
			ids = append(ids, val.ID)
		}
	}
	// reload corresponding nginx's container
	for _, id := range ids {
		args := fmt.Sprintf("exec %s nginx -s reload", id)
		res, err := execCommand("docker", args)
		if err != nil {
			fmt.Println(prefix + err.Error())
		} else {
			fmt.Println(res)
		}
	}
}

func execCommand(cmd string, args string) ([]string, error) {
	prefix := "[DnsManager][execCommand]"
	res := exec.Command(cmd, strings.Split(args, " ")...)
	stdout, err := res.StdoutPipe()
	res.Stderr = os.Stderr // read error from os
	err = res.Start()
	if err != nil {
		fmt.Println(prefix + err.Error())
		return nil, err
	}
	reader := bufio.NewReader(stdout)
	var results []string
	for {
		line, err2 := reader.ReadString('\n')
		if io.EOF == err2 {
			break
		} else {
			if err2 != nil {
				break
			}
		}
		results = append(results, line)
	}
	err = res.Wait()
	return results, err
}
