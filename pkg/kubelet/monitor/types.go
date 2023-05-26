package monitor

type StatsResponse struct {
	PodName        string  `json:"podName"`
	CPUUtilization float64 `json:"cpuUtilization"`
	MemoryUsage    float64 `json:"memoryUsage"`
	NetWorkIO      float64 `json:"netWorkIO,omitempty"`
}
