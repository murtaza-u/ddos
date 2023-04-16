package daemon

import (
	"context"
	"sync"

	log "github.com/murtaza-u/ddos/clog"
	"github.com/murtaza-u/ddos/direct"
	"github.com/murtaza-u/ddos/operate"
	pb "github.com/murtaza-u/ddos/proto/gen/go/daemon"
	"github.com/murtaza-u/ddos/store"
)

// DaemonSrv implements GRPC daemon server.
type DaemonSrv struct {
	Store     store.Storer
	Secret    string
	FileStore string

	pb.UnimplementedDaemonSvcServer
}

const BufferSize = 100

// Callback is a multi-operator, multi-method single channel of
// communication between the API server and daemon.
func (s *DaemonSrv) Callback(stream pb.DaemonSvc_CallbackServer) error {
	_id, _ := getID(stream.Context(), s.Secret)
	id := _id.String()

	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)

	defer func() {
		cancel()
		wg.Wait()
		log.Printf("[%s] Callback: Exiting...\n", id)
	}()

	var send = func(m any) error {
		return stream.SendMsg(m)
	}

	var recv = func() (direct.Request, error) {
		req, err := stream.Recv()
		return req, err
	}

	director := direct.NewSrvDirector(ctx, send, recv)

	wg.Add(1)
	go func() {
		err := director.Start()
		if err != nil {
			log.Printf("[%s] director.Start: %v\n", id, err)
		}
		cancel()
		wg.Done()
	}()

	out := director.Outlet()

Loop:
	for {
		var req direct.Request

		select {
		case <-ctx.Done():
			log.Printf("[%s] Callback: context canceled\n", id)
			break Loop
		case req = <-out:
			log.Printf(
				"[%s] Callback: received new request: %s | %s\n",
				id, req.GetIssuer(), req.GetMethod(),
			)
		}

		var op operate.Operator
		switch req.GetIssuer() {
		case pb.Issuer_ISSUER_DAEMON:
			op = operate.NewDaemonOperator(ctx, s.Store, id, req.GetId())
		case pb.Issuer_ISSUER_DDOS:
			op = operate.NewDDoSOperator(ctx, s.Store, id, req.GetId())
		default:
			continue Loop
		}

		switch req.GetMethod() {
		case pb.Method_METHOD_GET:
			go op.Get(director)
		case pb.Method_METHOD_PUT:
			go op.Put(director, req.(*pb.Request))
		case pb.Method_METHOD_WATCH:
			go op.Watch(director)
		case pb.Method_METHOD_APPEND:
			go op.Append(director, req.(*pb.Request))
		default:
			continue Loop
		}
	}

	return nil
}
