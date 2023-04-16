package mock

import (
	"context"
	"time"

	"github.com/murtaza-u/ddos/store"
)

type watcher struct {
	ctx    context.Context
	cancel context.CancelFunc
	evC    chan *store.Event
	errC   chan error
}

func (w watcher) Start() {
	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-t.C:
			ev := &store.Event{
				Val:     []byte{},
				Version: 0,
			}

			w.evC <- ev
		}
	}
}

func (w watcher) Stop() {
	w.cancel()
}

func (w watcher) EventChan() <-chan *store.Event {
	return w.evC
}

func (w watcher) ErrChan() <-chan error {
	return w.errC
}
