package daemon

import (
	"context"

	"github.com/murtaza-u/ddos/token"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

type intercept struct {
	secret string
}

func NewIntercept(secret string) *intercept {
	return &intercept{secret: secret}
}

func getID(ctx context.Context, secret string) (*uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"token not set",
		)
	}

	k := md.Get("token")[0]
	if k == "" {
		return nil, status.Errorf(
			codes.Unauthenticated,
			"token not set",
		)
	}

	t, err := token.Decrypt(k, secret, nil)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"failed to decrypt token: %v", err,
		)
	}

	id := new(uuid.UUID)
	err = t.Get("id", id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	return id, nil
}

// Stream intercepts the streaming gRPC request and validates the token.
func (i intercept) Stream(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {

	ctx := ss.Context()

	if info.FullMethod == "/DaemonSvc/Callback" {
		_, err := getID(ctx, i.secret)
		if err != nil {
			return err
		}
	}

	return handler(srv, &serverStream{
		ServerStream: ss,
		ctx:          ctx,
	})
}
