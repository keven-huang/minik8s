package informer

import (
	"minik8s/pkg/client/tool"
)

type EventHandler func(event tool.Event)

type Informer interface {
	AddEventHandler(etype tool.EventType, handler EventHandler)
	List() []tool.ListRes
	Run()
	Stop()
}

type informer struct {
	stop     bool
	resource string
	lw       tool.ListWatcher
	handlers map[tool.EventType]EventHandler
}

func (i *informer) AddEventHandler(etype tool.EventType, handler EventHandler) {
	i.handlers[etype] = handler
}

func (i *informer) List() []tool.ListRes {
	return i.lw.List(i.resource)
}

func (i *informer) Run() {
	i.stop = false
	watcher := i.lw.Watch(i.resource)
	for {
		select {
		case event := <-watcher.ResultChan():
			if i.stop {
				i.lw.Watch(i.resource).Stop()
				return
			}
			if handler, ok := i.handlers[event.Type]; ok {
				handler(event)
			}
		}
	}
}

func (i *informer) Stop() {
	i.stop = true
}

func NewInformer(resource string) Informer {
	return &informer{
		stop:     false,
		resource: resource,
		lw:       tool.NewListWatchFromClient(resource),
		handlers: make(map[tool.EventType]EventHandler, 4),
	}
}
