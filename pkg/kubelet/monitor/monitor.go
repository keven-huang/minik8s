package monitor

import (
	"encoding/json"
	"fmt"
	"io"
	"minik8s/cmd/kube-apiserver/app/apiconfig"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/tool"
	"minik8s/pkg/kubelet/dockerClient"
	"net/http"

	"github.com/docker/docker/api/types"
)

type Monitor struct {
	port int
	node *core.Node
}

func NewMonitor(port int, node *core.Node) *Monitor {
	return &Monitor{
		port: port,
		node: node,
	}
}

func (m *Monitor) Run() {
	http.HandleFunc("/stats", m.HandlerPodRequest)
	err := http.ListenAndServe(fmt.Sprintf(":%d", m.port), nil)
	if err != nil {
		panic(err)
	}
}

func readStats(reader io.Reader, v interface{}) error {
	dec := json.NewDecoder(reader)
	return dec.Decode(v)
}

func (m *Monitor) HandlerPodRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := tool.List(apiconfig.POD_PATH)
	resp := make([]StatsResponse, 0)
	for _, item := range res {
		var p core.Pod
		err := json.Unmarshal([]byte(item.Value), &p)
		if err != nil {
			fmt.Println("[kubelet][monitor] Unmarshal error:", err)
			continue
		}
		if p.Spec.NodeName == m.node.Name {
			stats, err := GetPodStats(&p)
			if err != nil {
				fmt.Println("[kubelet][monitor] GetPodStats error:", err)
				w.WriteHeader(http.StatusNotFound)
				return
			}
			resp = append(resp, stats)
		}
	}
	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("[kubelet][monitor] Marshal error:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetPodStats(pod *core.Pod) (StatsResponse, error) {
	var resp StatsResponse
	resp.PodName = pod.Name
	for _, container := range pod.Spec.Containers {
		stats, err := dockerClient.GetDockerStats(container.Name)
		if err != nil {
			return StatsResponse{}, err
		}
		var result types.StatsJSON
		err = readStats(stats.Body, &result)
		resp.CPUUtilization += getCPUPercent(&result)
		resp.MemoryUsage += getMemPercent(&result)
		fmt.Printf("CPU Usage: %f\n", resp.CPUUtilization)
		fmt.Printf("Memory Usage: %f\n", resp.CPUUtilization)
	}
	return resp, nil
}

func getCPUPercent(statsJson *types.StatsJSON) float64 {
	// cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	var preCPUUsage uint64
	var CPUUsage uint64
	preCPUUsage = 0
	CPUUsage = 0
	for _, core := range statsJson.CPUStats.CPUUsage.PercpuUsage {
		CPUUsage += core
	}
	for _, core := range statsJson.PreCPUStats.CPUUsage.PercpuUsage {
		preCPUUsage += core
	}
	systemUsage := statsJson.CPUStats.SystemUsage
	preSystemUsage := statsJson.PreCPUStats.SystemUsage

	deltaCPU := CPUUsage - preCPUUsage
	deltaSystem := systemUsage - preSystemUsage

	onlineCPU := statsJson.CPUStats.OnlineCPUs

	cpuPercent := (float64(deltaCPU) / float64(deltaSystem)) * float64(onlineCPU) * 100.0
	return cpuPercent
}

func getMemPercent(statsJson *types.StatsJSON) float64 {
	// MEM USAGE / LIMIT
	usage := statsJson.MemoryStats.Usage
	maxUsage := statsJson.MemoryStats.Limit
	percentage := float64(usage) / float64(maxUsage) * 100.0
	return percentage
}
