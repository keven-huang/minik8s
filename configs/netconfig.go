package configs

const NginxPath = "/root/nginx"   // host
const NginxPrefix = "/root/nginx" // nginx container prefix

const GatewayContainerPrefix = "GatewayContainer" // container prefix
const GatewayPodYamlPath = "./configs/gateway/gateway-pod.yaml"
const GatewayServiceYamlPath = "./configs/gateway/gateway-service.yaml"
const GatewayPodPrefix = "GatewayPod"
const GatewayServicePrefix = "GatewayService"

const CoreDnsPath = "../configs/dns/hostfile"
const CoreDnsPodYamlPath = "./configs/dns/coredns-pod.yaml"
const CoreDnsServiceYamlPath = "./configs/dns/coredns-service.yaml"
const CoreDNSServiceName = "coreDNS" // service name for coredns, that is unique
const CoreDNSPodName = "coreDNS"
