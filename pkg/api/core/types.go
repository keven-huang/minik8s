package core

type Container struct {
	// container's name
	Name  string `json:"name" yaml:"name"`
	Image string `json:"image" yaml:"image"`
	// init cmd when create this container,
	// 当command和entryPoint同时使用时，command作为entrypoint的参数
	Command    []string `json:"command" yaml:"command"`
	EntryPoint []string `json:"entryPoint" yaml:"entryPoint"`
	// ports
	Ports []Port `json:"ports" yaml:"ports"`
	// limit resource
	LimitResource Limit `json:"limitResource" yaml:"limitResource"`
	// TTY, if it is true and entry-point is 'sh', the container will not exit
	Tty bool `json:"tty" yaml:"tty"`
	// the mounted volumes in a pod
	VolumeMounts []VolumeMount `yaml:"volumeMounts" json:"volumeMounts"`
}

type VolumeMount struct {
	// volumeName
	Name string `yaml:"name" json:"name"`
	// mountpath
	MountPath string `yaml:"mountPath" json:"mountPath"`
}

type Port struct {
	Protocol   string // tcp, udp
	PortNumber string // number
}

type Limit struct {
	CPU    string `yaml:"cpu" json:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
}

type ContainerMeta struct {
	Name string
	Id   string
}
