package kube_service

import (
	"minik8s/pkg/api/core"
	"minik8s/pkg/client/informer"
	kube_service "minik8s/pkg/service"
	"sync"
	"time"
)

// 运行时service的配置
type RuntimeService struct {
	// 对应静态的service的配置
	ServiceConfig *kube_service.Service `yaml:"serviceConfig" json:"serviceConfig""`
	// 符合条件的pods
	Pods []core.Pod `json:"pods" yaml:"pods"`
	// 这些pods所在的ip
	// 可以通过endpoint+targetIP的方式来访问到这个pod的服务
	EndPoints []string `json:"endPoints" yaml:"endPoints"`
	// Informer是用于从api-server获取信息的
	Informer *informer.Informer
	// lock 用于同步
	lock sync.RWMutex
	// 定时器，更新Service最新信息
	timer *time.Ticker
	// 接受各个信息的channel
	eventChan chan string
	// 是否可以发送信号，如果为1表示可以
	ifSend bool
}

const TICK_EVENT = "tick"
