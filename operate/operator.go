package operate

import (
	"context"
	"errors"
	"sync"

	"github.com/murtaza-u/ddos/direct"
	pb "github.com/murtaza-u/ddos/proto/gen/go/daemon"
	"github.com/murtaza-u/ddos/store"
	"github.com/murtaza-u/ddos/types"
)

type Operator interface {
	Get(direct.SrvDirector) error
	Put(direct.SrvDirector, *pb.Request) error
	Append(direct.SrvDirector, *pb.Request) error
	Watch(direct.SrvDirector)
}

func NewDefaultOperator(
	ctx context.Context,
	s store.Storer,
	resource string,
	id string) Operator {

	return &defaultOpt{
		ctx:      ctx,
		store:    s,
		resource: resource,
		id:       id,
	}
}

type defaultOpt struct {
	ctx      context.Context
	store    store.Storer
	resource string
	id       string
}

func (o defaultOpt) Get(d direct.SrvDirector) error {
	out := &pb.Response{Id: o.id}

	v, ver, err := o.store.Get(o.resource)
	if err != nil {
		out.Err = &pb.Error{Error: err.Error()}
		return d.Direct(out)
	}

	out.Err = nil
	out.Resource = &pb.Resource{
		Version:  ver,
		Manifest: v,
	}

	return d.Direct(out)
}

func (o defaultOpt) Put(d direct.SrvDirector, in *pb.Request) error {
	out := &pb.Response{Id: o.id}

	v, ver, err := o.store.Get(o.resource)
	found := !errors.Is(err, store.ErrKeyNotFound)

	if err != nil && found {
		out.Err = &pb.Error{Error: err.Error()}
		return d.Direct(out)
	}

	if found && ver > in.Resource.GetVersion() {
		out.Err = &pb.Error{Error: types.ErrKeyModified.Error()}
		out.Resource = &pb.Resource{
			Version:  ver,
			Manifest: v,
		}
		return d.Direct(out)
	}

	ver, err = o.store.Put(o.resource, in.Resource.GetManifest())
	if err != nil {
		out.Err = &pb.Error{Error: err.Error()}
		return d.Direct(out)
	}

	out.Err = nil
	out.Resource = &pb.Resource{
		Version:  ver,
		Manifest: in.Resource.GetManifest(),
	}

	return d.Direct(out)
}

func (o defaultOpt) Append(d direct.SrvDirector, _ *pb.Request) error {
	out := &pb.Response{Id: o.id}
	out.Err = &pb.Error{
		Error: types.ErrUnsupportedMethod.Error(),
	}
	return d.Direct(out)
}

func (o defaultOpt) Watch(d direct.SrvDirector) {
	ctx, cancel := context.WithCancel(o.ctx)
	ch := make(chan direct.IdGetter, 100)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		d.WatcherDirect(ch)
		cancel()
		wg.Done()
	}()

	w := o.store.Watcher(ctx, o.resource)
	defer w.Stop()

	wg.Add(1)
	go func() {
		w.Start()
		cancel()
		wg.Done()
	}()

Loop:
	for {
		out := &pb.Response{Id: o.id}
		select {
		case <-ctx.Done():
			break Loop
		case ev := <-w.EventChan():
			switch ev.Type {
			case store.EventAdded:
				out.Event = pb.Event_EVENT_ADDED
			case store.EventModified:
				out.Event = pb.Event_EVENT_MODIFIED
			case store.EventDeleted:
				out.Event = pb.Event_EVENT_DELETED
			default:
				continue Loop
			}

			out.Resource = &pb.Resource{
				Version:  ev.Version,
				Manifest: ev.Val,
			}
			ch <- out
		case err := <-w.ErrChan():
			out.Err = &pb.Error{Error: err.Error()}
			ch <- out
			break Loop
		}
	}

	wg.Wait()
}
