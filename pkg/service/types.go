package service

import metav1 "minik8s/pkg/apis/meta/v1"

type Service struct {
	ServiceMeta metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	ServiceSpec Spec              `json:"spec" yaml:"spec"`
}

type Spec struct {
	// service的name
	Name string `json:"name" yaml:"name"`
	// selector 用于筛选pod
	Selector map[string]string `json:"selector" yaml:"selector"`
	// service的类型，可能是ClusterIP, NodePort ...
	Type Type `json:"type" yaml:"type"`
	// service的端口信息
	Ports []Port `json:"ports" yaml:"ports"`
	// 用户指定的clusterIP, 对type有要求
	ClusterIP string `json:"clusterIP" yaml:"clusterIP"`
	// 状态
	Status ServiceStatus `json:"status" yaml:"status"`
}

type ServiceStatus struct {
	Err   error  `yaml:"err" json:"err"`
	Phase string `yaml:"phase" json:"phase"` // creating, running, error
}

const (
	ServiceRunningPhase  string = "running"
	ServiceCreatingPhase string = "creating"
	ServiceErrorPhase    string = "error"
)

type Type string

const (
	ClusterIPTYPE Type = "CluterIP"
	NodePortTYPE  Type = "NodePort"
)

type Port struct {
	// tcp, udp
	Protocol string `yaml:"protocol" json:"protocol"`
	// 端口，是service对外的端口
	Port string `json:"port" yaml:"port"`
	// 这个service使用的pod对外的端口
	TargetPort string `yaml:"targetPort" json:"targetPort"`
}

type Status struct {
	// Delete, running, pending
	Phase string `yaml:"phase" json:"phase"`
}
