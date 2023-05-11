package kube_proxy

// serviceIP chain
type SipChain struct {
	// 名字
	Name     string
	Protocol string
	Table    string
	// 指向的pod的IP
	PodIP            string
	PodName          string
	DNatRule         *DNatRule
	RoundRabinNumber int
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
}
