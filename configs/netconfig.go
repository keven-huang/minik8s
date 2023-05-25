package configs

const NginxPath = "/root/nginx"   // host
const NginxPrefix = "/root/nginx" // nginx container prefix

const GatewayContainerPrefix = "GatewayContainer" // container prefix
const GatewayPodYamlPath = "/home/minik8s/configs/gateway/gateway-pod.yaml"
const GatewayServiceYamlPath = "/home/minik8s/configs/gateway/gateway-service.yaml"
const GatewayPodPrefix = "GatewayPod"
const GatewayServicePrefix = "GatewayService"

const CoreDnsPath = "/home/minik8s/configs/dns/hostfile"
const CoreDnsPodYamlPath = "/home/minik8s/configs/dns/coredns-pod.yaml"
const CoreDnsServiceYamlPath = "/home/minik8s/configs/dns/coredns-service.yaml"
const CoreDNSServiceName = "coreDNS" // service name for coredns, that is unique
const CoreDNSPodName = "coreDNS"
