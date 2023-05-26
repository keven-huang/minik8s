package monitor

type StatsResponse struct {
	CPUUtilization float64 `json:"cpuUtilization,omitempty"`
	MemoryUsage    float64 `json:"memoryUsage,omitempty"`
	NetWorkIO      float64 `json:"netWorkIO,omitempty"`
}
