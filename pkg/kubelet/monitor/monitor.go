package monitor

import (
	"encoding/json"
	"fmt"
	"io"
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/tool"
	"minik8s/pkg/kubelet/dockerClient"
	"net/http"

	"github.com/docker/docker/api/types"
)

type Monitor struct {
	port int
}

func NewMonitor(port int) *Monitor {
	return &Monitor{
		port: port,
	}
}

func (m *Monitor) Run() {
	http.HandleFunc("/stats", HandlerPodRequest)
	err := http.ListenAndServe(fmt.Sprintf(":%d", m.port), nil)
	if err != nil {
		panic(err)
	}
}

func readStats(reader io.Reader, v interface{}) error {
	dec := json.NewDecoder(reader)
	return dec.Decode(v)
}

func HandlerPodRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryParams := r.URL.Query()
	podName := queryParams.Get("podName")

	pod, err := tool.GetPod(podName)
	if err != nil {
		fmt.Println("[kubelet][monitor] GetPod error:", err)
		w.WriteHeader(http.StatusNotFound) // 返回404状态码
		return
	}
	resp, err := GetPodStats(pod)
	if err != nil {
		fmt.Println("[kubelet][monitor] GetPodStats error:", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetPodStats(pod *core.Pod) (StatsResponse, error) {
	var resp StatsResponse
	for _, container := range pod.Spec.Containers {
		stats, err := dockerClient.GetDockerStats(container.Name)
		if err != nil {
			return StatsResponse{}, err
		}
		var result types.StatsJSON
		err = readStats(stats.Body, &result)
		resp.CPUUtilization = getCPUPercent(&result)
		resp.MemoryUsage = getMemPercent(&result)
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
