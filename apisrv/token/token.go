package token

import (
	"context"
	"fmt"
	"time"

	"github.com/murtaza-u/ddos/manifest"
	"github.com/murtaza-u/ddos/manifest/daemon"
	pb "github.com/murtaza-u/ddos/proto/gen/go/token"
	"github.com/murtaza-u/ddos/store"
	"github.com/murtaza-u/ddos/token"
	"github.com/murtaza-u/ddos/types"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TokenSrv implements GRPC token server.
type TokenSrv struct {
	Secret string
	Store  store.Storer

	pb.UnimplementedTokenSvcServer
}

// Register generates an new UID, stores the host information in the
// database and returns a uniquely identifiable token.
func (s *TokenSrv) Register(ctx context.Context, info *pb.HostInfo) (*pb.Token, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, "failed to create new UID",
		)
	}

	p := token.Params{
		Issuer:   "/TokenSvc",
		Audience: "/DaemonSvc",
		Body: map[string]any{
			"role": token.RoleDaemon,
			"id":   id,
		},
	}

	t, err := token.New(p)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, "failed to create new token: %v",
			err,
		)
	}

	enc, err := t.Encrypt(s.Secret, nil)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, "failed to encrypt token: %v", err,
		)
	}

	k := fmt.Sprintf("/registry/%s/daemon", id.String())
	v := &daemon.Manifest{
		Metadata: daemon.Metadata{
			Name: id.String(),
			Info: daemon.Info{
				Hostname: info.GetHostname(),
				Username: info.GetUsername(),
				OS:       info.GetOs(),
				Platform: info.GetPlatform(),
				Cores:    info.GetCpucount(),
				Memory:   info.GetMemtotal(),
			},
			LastPing: time.Now().UTC(),
		},
	}

	out, err := manifest.Marshal(v)
	if err != nil {
		return nil, types.ErrMarshaling
	}

	_, err = s.Store.Put(k, out)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, "failed to store host info: %v", err,
		)
	}

	return &pb.Token{T: enc}, nil
}
