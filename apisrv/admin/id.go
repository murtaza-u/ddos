package admin

import (
	"context"
	"strings"

	pb "github.com/murtaza-u/ddos/proto/gen/go/admin"
	"github.com/murtaza-u/ddos/store"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func idsFromStore(store store.Storer) ([]string, error) {
	keys, err := store.GetKeysWithPrefix("/registry/")
	if err != nil {
		return nil, grpc.Errorf(
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

	return ids, nil
}

func (a *AdminSrv) GetIds(ctx context.Context, _ *pb.Void) (*pb.Ids, error) {
	ids, err := idsFromStore(a.Store)
	if err != nil {
		return nil, err
	}
	return &pb.Ids{Ids: ids}, nil
}
