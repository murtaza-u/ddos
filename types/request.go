package types

import pb "github.com/murtaza-u/ddos/proto/gen/go/daemon"

type Request struct {
	*pb.Request
}

func (r *Request) SetId(id string) {
	r.Id = id
}
