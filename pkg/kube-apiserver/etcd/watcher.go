package etcd

import clientv3 "go.etcd.io/etcd/client/v3"

type EventType string

const (
	Added    EventType = "ADDED"
	Modified EventType = "MODIFIED"
	Deleted  EventType = "DELETED"
	Error    EventType = "ERROR"
)

type Event struct {
	Type EventType
	Key  string
	Val  string
}

type ListRes struct {
	ResourceVersion string
	Key             string
	Value           string
}

// help function
func getType(event *clientv3.Event) EventType {
	switch event.Type {
	case clientv3.EventTypePut:
		if event.IsCreate() {
			return Added
		} else if event.IsModify() {
			return Modified
		} else {
			return Error
		}
	case clientv3.EventTypeDelete:
		return Deleted
	default:
		return Error
	}
}
