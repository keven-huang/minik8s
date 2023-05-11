package kube_proxy

func (rule *DNatRule) toSpec() {
	// table
	rule.Spec = []string{"-t", rule.Table}
	// portocol
	rule.Spec = append(rule.Spec, "-p", rule.Protocol)
	// source
	rule.Spec = append(rule.Spec, "-s", rule.SourceIP, "--sport", rule.SourcePort)
	// dest
	rule.Spec = append(rule.Spec, "-d", rule.DestIP, "--dport", rule.DestPort)
	// action, 指做 入方向的地址转换
	rule.Spec = append(rule.Spec, "-j", "DNAT")

}

func NewDNatRule(PodIP string,
	port string,
	SourceIP *string,
	SourcePort *string,
	protocol string,
	table *string) *DNatRule {
	res := &DNatRule{}
	res.DestIP = PodIP
	res.DestPort = port
	if SourceIP != nil {
		res.SourceIP = *SourceIP
	} else {
		res.SourceIP = "0"
	}
	if SourcePort != nil {
		res.SourcePort = *SourcePort
	} else {
		res.SourcePort = "0"
	}
	res.Protocol = protocol
	if table != nil {
		res.Table = *table
	} else {
		res.Table = "nat"
	}
	res.toSpec()
	return res
}
