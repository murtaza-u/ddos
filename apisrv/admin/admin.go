package admin

import (
	pb "github.com/murtaza-u/ddos/proto/gen/go/admin"
	"github.com/murtaza-u/ddos/store"
)

// AdminSrv implements GRPC admin server.
type AdminSrv struct {
	Store     store.Storer
	FileStore string

	pb.UnimplementedAdminSvcServer
}
