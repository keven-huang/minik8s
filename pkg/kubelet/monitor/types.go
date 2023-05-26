package monitor

type StatsResponse struct {
	CPUUtilization float64 `json:"cpuUtilization"`
	MemoryUsage    float64 `json:"memoryUsage"`
	NetWorkIO      float64 `json:"netWorkIO,omitempty"`
}
