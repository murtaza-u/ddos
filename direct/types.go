package direct

import pb "github.com/murtaza-u/ddos/proto/gen/go/daemon"

const bufsize = 100

type Send func(any) error
type SrvRecv func() (Request, error)
type DaemonRecv func() (Response, error)

type IdGetter interface {
	GetId() string
}

type idSetter interface {
	SetId(string)
}

type IdGetSetter interface {
	IdGetter
	idSetter
}

type issuer interface {
	GetIssuer() pb.Issuer
}

type methoder interface {
	GetMethod() pb.Method
}

type Request interface {
	IdGetter
	issuer
	methoder
}

type Response interface {
	IdGetter
}
