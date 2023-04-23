package core

import (
	"minik8s/pkg/kubelet"
)

type Container struct {
	// container's name
	Name  string `json:"name" yaml:"name"`
	Image string `json:"image" yaml:"image"`
	// init cmd when create this container,
	// 当command和entryPoint同时使用时，command作为entrypoint的参数
	Command    []string `json:"command" yaml:"command"`
	EntryPoint []string `json:"entryPoint" yaml:"entryPoint"`
	// ports
	Ports []kubelet.Port `json:"ports" yaml:"ports"`
	// limit resource
	LimitResource Limit `json:"limitResource" yaml:"limitResource"`
}

type Limit struct {
	CPU    string `yaml:"cpu" json:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
}

type ContainerMeta struct {
	Name string
	Id   string
}
