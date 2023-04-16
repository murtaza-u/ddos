package operate

import (
	"context"
	"fmt"
	"time"

	"github.com/murtaza-u/ddos/direct"
	"github.com/murtaza-u/ddos/manifest"
	pb "github.com/murtaza-u/ddos/proto/gen/go/daemon"
	"github.com/murtaza-u/ddos/store"
)

type daemonOpt struct {
	ctx      context.Context
	store    store.Storer
	resource string
	id       string
}

func NewDaemonOperator(
	ctx context.Context,
	s store.Storer,
	name string,
	id string) Operator {

	resource := fmt.Sprintf("/registry/%s/daemon", name)
	return &daemonOpt{
		ctx:      ctx,
		store:    s,
		resource: resource,
		id:       id,
	}
}

func (o daemonOpt) Get(d direct.SrvDirector) error {
	return NewDefaultOperator(o.ctx, o.store, o.resource, o.id).Get(d)
}

func (o daemonOpt) Put(d direct.SrvDirector, in *pb.Request) error {
	return NewDefaultOperator(o.ctx, o.store, o.resource, o.id).Put(d, in)
}

func (o daemonOpt) Append(d direct.SrvDirector, in *pb.Request) error {
	out := &pb.Response{Id: o.id}
	v, ver, err := o.store.Get(o.resource)
	if err != nil {
		out.Err = &pb.Error{Error: err.Error()}
		return d.Direct(out)
	}

	m := manifest.NewDaemonManifest()
	err = manifest.Unmarshal(v, m)
	if err != nil {
		out.Err = &pb.Error{Error: err.Error()}
		return d.Direct(out)
	}

	m.Metadata.LastPing = time.Now().UTC()
	v, err = manifest.Marshal(m)
	if err != nil {
		out.Err = &pb.Error{Error: err.Error()}
		return d.Direct(out)
	}

	ver, err = o.store.Put(o.resource, v)
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

func (o daemonOpt) Watch(d direct.SrvDirector) {
	NewDefaultOperator(o.ctx, o.store, o.resource, o.id).Watch(d)
}
