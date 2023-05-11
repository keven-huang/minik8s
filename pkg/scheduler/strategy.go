package scheduler

import "minik8s/pkg/api/core"

var global_cnt int = 0

func roundrobin_strategy(nodes []core.Node) string {
	global_cnt = (global_cnt + 1) % len(nodes)
	return nodes[global_cnt].Name
}
