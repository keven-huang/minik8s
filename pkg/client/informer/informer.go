package informer

import (
	"minik8s/pkg/client/tool"
)

type EventHandler func(event tool.Event)

type Informer interface {
	AddEventHandler(etype tool.EventType, handler EventHandler)
	Run()
}

type informer struct {
	resource string
	lw       tool.ListWatcher
	handlers map[tool.EventType]EventHandler
}

func (i *informer) AddEventHandler(etype tool.EventType, handler EventHandler) {
	i.handlers[etype] = handler
}

func (i *informer) Run() {
	watcher := i.lw.Watch(i.resource)
	for {
		select {
		case event := <-watcher.ResultChan():
			if handler, ok := i.handlers[event.Type]; ok {
				handler(event)
			}
		}
	}
	// i.lw.Watch(i.resource).Stop()
}

func NewInformer(resource string) Informer {
	return &informer{
		resource: resource,
		lw:       tool.NewListWatchFromClient(resource),
		handlers: make(map[tool.EventType]EventHandler, 4),
	}
}
