package direct

import (
	"context"

	"github.com/murtaza-u/ddos/types"
)

type SrvDirector interface {
	Direct(IdGetter) error
	WatcherDirect(<-chan IdGetter) error
	Outlet() <-chan Request
	Start() error
}

func NewSrvDirector(ctx context.Context, s Send, r SrvRecv) SrvDirector {
	ctx, cancel := context.WithCancel(ctx)

	return &sdirect{
		ctx:    ctx,
		cancel: cancel,
		out:    make(chan Request, bufsize),
		send:   s,
		recv:   r,
	}
}

type sdirect struct {
	ctx    context.Context
	cancel context.CancelFunc
	out    chan Request
	send   Send
	recv   SrvRecv
}

func (d *sdirect) Outlet() <-chan Request {
	return d.out
}

func (d *sdirect) Direct(m IdGetter) error {
	id := m.GetId()
	if id == "" {
		return types.ErrIdNotSet
	}

	err := d.send(m)
	if err != nil {
		d.cancel()
		return err
	}

	return nil
}

func (d *sdirect) WatcherDirect(ch <-chan IdGetter) error {
	for {
		select {
		case <-d.ctx.Done():
			return nil
		case m := <-ch:
			id := m.GetId()
			if id == "" {
				return types.ErrIdNotSet
			}

			err := d.send(m)
			if err != nil {
				d.cancel()
				return err
			}
		}
	}
}

func (d *sdirect) Start() error {
	for {
		m, err := d.recv()
		if err != nil {
			d.cancel()
			return err
		}

		d.out <- m
	}
}
