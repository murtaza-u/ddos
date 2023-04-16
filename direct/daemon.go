package direct

import (
	"context"
	"log"
	"sync"

	"github.com/murtaza-u/ddos/types"

	"github.com/google/uuid"
)

type DaemonDirector interface {
	Start() error
	Direct(IdGetSetter) (any, error)
	WatcherDirect(IdGetSetter) (<-chan IdGetter, error)
}

func NewDaemonDirector(ctx context.Context, s Send, r DaemonRecv) DaemonDirector {
	ctx, cancel := context.WithCancel(ctx)

	return &ddirect{
		ctx:    ctx,
		cancel: cancel,
		send:   s,
		recv:   r,
		req:    make(chan IdGetter, bufsize),
		Mutex:  sync.Mutex{},
		table:  make(map[string]chan IdGetter),
	}
}

type ddirect struct {
	ctx    context.Context
	cancel context.CancelFunc

	send Send
	recv DaemonRecv
	req  chan IdGetter

	sync.Mutex
	table map[string]chan IdGetter
}

func (d *ddirect) Direct(m IdGetSetter) (any, error) {
	id := uuid.NewString()
	m.SetId(id)

	ch := d.addEntry(id)
	defer d.delEntry(id)

	err := d.send(m)
	if err != nil {
		d.cancel()
		return nil, err
	}

	select {
	case <-d.ctx.Done():
		return nil, types.ErrCanceled
	case m := <-ch:
		return m, nil
	}
}

func (d *ddirect) WatcherDirect(m IdGetSetter) (<-chan IdGetter, error) {
	id := uuid.NewString()
	m.SetId(id)

	in := d.addEntry(id)
	out := make(chan IdGetter, bufsize)

	go func(in <-chan IdGetter, out chan<- IdGetter) {
		defer d.delEntry(id)

		for {
			select {
			case <-d.ctx.Done():
				return
			case m := <-in:
				out <- m
			}
		}
	}(in, out)

	err := d.send(m)
	if err != nil {
		d.cancel()
		return nil, err
	}

	return out, nil
}

func (d *ddirect) Start() error {
	go d.push()

	for {
		select {
		case <-d.ctx.Done():
			return types.ErrCanceled
		case m := <-d.req:
			id := m.GetId()
			if id == "" {
				continue
			}

			d.Lock()
			ch := d.table[id]
			d.Unlock()

			if ch == nil {
				log.Println(
					"[DaemonDirector] table entry does not exist. Skipping",
				)
				continue
			}

			ch <- m
		}
	}
}

func (d *ddirect) push() error {
	for {
		m, err := d.recv()
		if err != nil {
			log.Println(err)
			d.cancel()
			return err
		}

		d.req <- m
	}
}

func (d *ddirect) addEntry(id string) <-chan IdGetter {
	d.Lock()
	defer d.Unlock()

	d.table[id] = make(chan IdGetter)
	return d.table[id]
}

func (d *ddirect) delEntry(k string) {
	d.Lock()
	defer d.Unlock()

	delete(d.table, k)
}
