package kube_proxy

import (
	"bufio"
	"fmt"
	"io"
	"minik8s/configs"
	"minik8s/pkg/api/core"
	"minik8s/pkg/kubelet/dockerClient"
	"os"
	"os/exec"
	"strings"
)

func (proxy *KubeProxy) writeNginx(dns *core.DNS) {
	prefix := "[DNSManager][WriteNginx]"
	fmt.Println(prefix + "in")
	var data []string
	data = append(data, "user  root;")
	data = append(data, "worker_processes  1;")
	data = append(data, "events {")
	data = append(data, "    worker_connections 1024;")
	data = append(data, "}")
	data = append(data, "http {")
	data = append(data, "    server {", "        listen 80;")
	//data = append(data, fmt.Sprintf("        server_name %s;", dns.Spec.Host)) // it can be deleted maybe
	data = append(data, proxy.generateConfig(dns)...)
	data = append(data, "    }")
	data = append(data, "}")
	file, err := os.OpenFile(configs.NginxPrefix+"/"+dns.Metadata.Name+"/"+"nginx.conf", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
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
func (proxy *KubeProxy) generateConfig(dns *core.DNS) []string {
	var res []string
	for _, v := range dns.Spec.Paths {
		res = append(res, fmt.Sprintf("        location %s {", v.Name))
		if v.Ip != "" {
			res = append(res, fmt.Sprintf("            proxy_pass http://%s:%s/;", v.Ip, v.Port))
		}
		res = append(res, "        }")
	}
	return res
}

func (proxy *KubeProxy) mkDir(Dns string) {
	prefix := "[DNSManager][mkDir]"
	fmt.Println(prefix + "in")
	args := fmt.Sprintf("%s", configs.NginxPrefix+"/"+Dns)
	res, err := execCommand("mkdir", args)
	if err != nil {
		fmt.Println(prefix + err.Error())
	} else {
		fmt.Println(prefix + "resValue:")
		fmt.Println(res)
	}
}

func (proxy *KubeProxy) deleteDir(Dns string) {
	prefix := "[DNSManager][deleteDir]"
	fmt.Println(prefix + "in")
	args := fmt.Sprintf("-rf %s", configs.NginxPrefix+"/"+Dns)
	res, err := execCommand("rm", args)
	if err != nil {
		fmt.Println(prefix + err.Error())
	} else {
		fmt.Println(res)
	}
}

// update dns entry based on key2DNs
// mapping from domainName -> gateway
// the gateway is a nginx
func (proxy *KubeProxy) writeCoreDNS() {
	prefix := "[DNSManager][WriteCoreDNS]"
	wd, err := os.Getwd()
	fmt.Println(prefix + "in:" + wd)
	file, err := os.OpenFile(configs.CoreDnsPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println(prefix + err.Error())
		return
	}
	writer := bufio.NewWriter(file)
	for _, v := range proxy.Key2Dns {
		if v.Status == core.ServiceCreatedStatus {
			cur := v.Spec.GatewayIp + " " + v.Spec.Host
			_, err := fmt.Fprintln(writer, cur)
			if err != nil {
				fmt.Println(prefix + err.Error())
			} else {
				fmt.Println(prefix + " success add mappings" + v.Spec.GatewayIp + ":" + v.Spec.Host)
			}
		}
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println(prefix + err.Error())
	}
	return
}

func (proxy *KubeProxy) reloadCoreDNS() {
	prefix := "[DNSManager][reloadCoreDNS]"
	fmt.Println(prefix + "in")
	cons, err := dockerClient.GetAllContainers()
	if err != nil {
		fmt.Println(prefix + err.Error())
		return
	}
	var ids []string
	for _, val := range cons {
		// TODO names
		if strings.Contains(val.Names[0], configs.CoreDNSPodName) {
			ids = append(ids, val.ID)
		}
	}
	// reload corresponding container
	for _, id := range ids {
		fmt.Println(prefix + "reloading..")
		args := fmt.Sprintf("exec -d %s /coredns -conf /etc/coredns/Corefile", id)
		res, err := execCommand("docker", args)
		if err != nil {
			fmt.Println(prefix + err.Error())
		} else {
			fmt.Println(prefix + "return value:")
			fmt.Println(res)
		}
	}
}

func (proxy *KubeProxy) reloadNginx(dns *core.DNS) {
	prefix := "[DNSManager][reloadNginx]"
	fmt.Println(prefix + "in")
	cons, err := dockerClient.GetAllContainers()
	if err != nil {
		fmt.Println(prefix + err.Error())
		return
	}
	var ids []string
	for _, val := range cons {
		// TODO names
		if strings.Contains(val.Names[0], configs.GatewayContainerPrefix+dns.Metadata.Name) {
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
