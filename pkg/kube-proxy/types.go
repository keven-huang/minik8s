package kube_proxy

// PodChain pod chain, 一条链对应一个pod
type PodChain struct {
	// 名字
	Name        string
	Protocol    string
	Table       string
	FatherChain string
	// 对应的pod的信息
	Pod      PodInfo
	DNatRule *DNatRule
	// 每n个包执行一次规则
	RoundRabinNumber int
	// 格式化信息
	Spec []string
}

const PodChainPrefix = "Pod"
const SvcChainPrefix = "Svc"

// SvcChain service chain, 对应某个service
type SvcChain struct {
	Name        string
	Table       string
	FatherChain string
	Protocol    string
	// svc对外暴露的集群ip
	ClusterIp   string
	ClusterPort string
	// PodName->PodChain
	Name2Chain map[string]*PodChain
	Spec       []string
}

type PodInfo struct {
	IP   string
	Name string
	Port string
}

type DNatRule struct {
	// 目标IP
	DestIP string
	// Dest port
	DestPort string
	// Source Ip
	SourceIP string
	// Source Port
	SourcePort string
	// Protocol
	Protocol string
	// Spec format for linux
	Spec []string
	// Table Name, default is Nat
	Table string
	// father chain, 父链
	FatherChain string
}
