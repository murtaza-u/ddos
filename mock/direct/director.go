package mock

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/murtaza-u/ddos/direct"
	pb "github.com/murtaza-u/ddos/proto/gen/go/daemon"
)

type counter struct {
	sync.Mutex
	i int
}

func Reset() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.i = 0
}

var c = &counter{
	i:     0,
	Mutex: sync.Mutex{},
}

type Request struct {
	Id string
}

func (r Request) GetId() string {
	return r.Id
}

func (r *Request) SetId(id string) {
	r.Id = id
}

func (Request) GetIssuer() pb.Issuer {
	return pb.Issuer_ISSUER_UNSPECIFIED
}

func (Request) GetMethod() pb.Method {
	return pb.Method_METHOD_UNSPECIFIED
}

func recv() (direct.Request, error) {
	time.Sleep(time.Second)

	c.Lock()
	defer c.Unlock()

	if c.i >= 5 {
		return nil, errors.New("fake client closed the stream")
	}

	c.i++
	log.Println("receive request received")
	return &Request{Id: "XXX"}, nil
}

func send(m any) error {
	time.Sleep(time.Second)

	c.Lock()
	defer c.Unlock()

	if c.i >= 5 {
		return errors.New("fake client closed the stream")
	}

	c.i++
	log.Println("send request receieved")
	return nil
}

func NewSrvDirector(ctx context.Context) direct.SrvDirector {
	return direct.NewSrvDirector(ctx, send, recv)
}
