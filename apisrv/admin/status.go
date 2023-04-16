package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/murtaza-u/ddos/manifest"
	"github.com/murtaza-u/ddos/manifest/daemon"
	pb "github.com/murtaza-u/ddos/proto/gen/go/admin"
	"github.com/murtaza-u/ddos/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Status returns an array of daemon manifests.
func (a *AdminSrv) Status(ctx context.Context, _ *pb.Void) (*pb.Byts, error) {
	keys, err := a.Store.GetKeysWithPrefix("/registry/")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, "failed to get keys from db: %v", err,
		)
	}

	var ids []string

Outer:
	for _, k := range keys {
		id := strings.Split(k, "/")[2]
		for _, i := range ids {
			if i == id {
				continue Outer
			}
		}

		ids = append(ids, id)
	}

	out := manifest.NewStatusManifest()
	out.Daemons = []daemon.Manifest{}

	for _, id := range ids {
		k := fmt.Sprintf("/registry/%s/daemon", id)
		v, _, err := a.Store.Get(k)
		if err != nil {
			if err == types.ErrKeyNotFound {
				continue
			}

			return nil, status.Errorf(
				codes.Internal, err.Error(),
			)
		}

		m := new(daemon.Manifest)
		err = manifest.Unmarshal([]byte(v), m)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal, types.ErrUnmarshaling.Error(),
			)
		}

		out.Daemons = append(out.Daemons, *m)
	}

	byts, err := manifest.Marshal(out)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, types.ErrMarshaling.Error(),
		)
	}

	return &pb.Byts{
		B: byts,
	}, nil
}
