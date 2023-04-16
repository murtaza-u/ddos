package store

import clientv3 "go.etcd.io/etcd/client/v3"

// Event is a database agnostic event object.
type Event struct {
	Type    EventType
	Val     []byte
	Version int64
}

type EventType string

const (
	EventAdded    EventType = "ADDED"
	EventModified EventType = "MODIFIED"
	EventDeleted  EventType = "DELETED"
)

// parseEvent converts *clientv3.Event to a database agnostic event.
func parseEvent(ev *clientv3.Event) *Event {
	event := new(Event)
	event.Val = ev.Kv.Value
	event.Version = ev.Kv.Version

	if ev.IsCreate() {
		event.Type = EventAdded
	}

	if ev.Type == clientv3.EventTypeDelete {
		event.Type = EventDeleted
	}

	if ev.Type == clientv3.EventTypePut {
		event.Type = EventModified
	}

	return event
}
