package informer

import (
	"minik8s/pkg/client/tool"
	"minik8s/pkg/kube-apiserver/etcd"
)

type EventHandler func(event tool.Event)

type Informer interface {
	AddEventHandler(handler EventHandler)
	Run()
}

type informer struct {
	resource   string
	lw         tool.ListWatcher
	informType etcd.EventType
	handler    EventHandler
}

func (i *informer) AddEventHandler(handler EventHandler) {
	i.handler = handler
}

func (i *informer) Run() {
	watcher := i.lw.Watch(i.resource)
	for {
		select {
		case event := <-watcher.ResultChan():
			if event.Type == i.informType {
				i.handler(event)
			}
		}
	}
	// i.lw.Watch(i.resource).Stop()
}

func NewInformer(resource string) Informer {
	return &informer{
		resource: resource,
		lw:       tool.NewListWatchFromClient(resource),
	}
}
