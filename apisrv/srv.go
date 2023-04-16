package apisrv

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/murtaza-u/ddos/apisrv/admin"
	"github.com/murtaza-u/ddos/apisrv/daemon"
	"github.com/murtaza-u/ddos/apisrv/token"
	"github.com/murtaza-u/ddos/conf"
	pbadm "github.com/murtaza-u/ddos/proto/gen/go/admin"
	pbd "github.com/murtaza-u/ddos/proto/gen/go/daemon"
	pbtkn "github.com/murtaza-u/ddos/proto/gen/go/token"
	"github.com/murtaza-u/ddos/store"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// Start starts the GRPC API server.
func Start(c *conf.C) error {
	err := c.Validate()
	if err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	addr := fmt.Sprintf(":%d", c.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	store, err := store.New(c.Endpoints)
	if err != nil {
		return fmt.Errorf("failed to access store: %w", err)
	}

	i := daemon.NewIntercept(c.Secret)

	var s *grpc.Server
	if c.Tls {
		creds, err := loadTLSCreds(c)
		if err != nil {
			return fmt.Errorf("failed to load tls credentials: %w", err)
		}

		s = grpc.NewServer(
			grpc.Creds(creds), grpc.StreamInterceptor(i.Stream),
		)
	} else {
		s = grpc.NewServer(grpc.StreamInterceptor(i.Stream))
	}

	pbtkn.RegisterTokenSvcServer(s, &token.TokenSrv{
		Secret: c.Secret,
		Store:  store,
	})
	pbd.RegisterDaemonSvcServer(s, &daemon.DaemonSrv{
		Store:     store,
		Secret:    c.Secret,
		FileStore: c.FileStore,
	})
	pbadm.RegisterAdminSvcServer(s, &admin.AdminSrv{
		Store:     store,
		FileStore: c.FileStore,
	})

	if c.Reflect {
		reflection.Register(s)
	}

	return s.Serve(ln)
}

func loadTLSCreds(c *conf.C) (credentials.TransportCredentials, error) {
	// load cert of CA
	caCert, err := os.ReadFile(filepath.Join(c.Certs, "ca-cert.pem"))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read %s/ca-cert.pem: %w", c.Certs, err,
		)
	}

	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(caCert)
	if !ok {
		return nil, fmt.Errorf(
			"failed to add CA's certificate to pool",
		)
	}

	// load apisrv's cert and priv key
	cert, err := tls.LoadX509KeyPair(
		filepath.Join(c.Certs, "apisrv-cert.pem"),
		filepath.Join(c.Certs, "apisrv-key.pem"),
	)
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
	}

	return credentials.NewTLS(cfg), nil
}
