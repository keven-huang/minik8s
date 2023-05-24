package kube_proxy

const NginxPath = "/root/nginx"                   // host
const NginxPrefix = "/root/nginx"                 // nginx container prefix
const GatewayContainerPrefix = "GatewayContainer" // container prefix
const CoreDnsPath = "/root/coredns/hostsfile"

const CoreDnsPodYamlPath = "/home/minik8s/configs/dns/coredns-pod.yaml"
const CoreDnsServiceYamlPath = "/home/minik8s/configs/dns/coredns-service.yaml"

const CoreDNSServiceName = "coreDNS" // service name for coredns, that is unique
