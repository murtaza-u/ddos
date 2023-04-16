package manifest

import (
	"github.com/murtaza-u/ddos/manifest/daemon"
	"github.com/murtaza-u/ddos/manifest/ddos"
	"github.com/murtaza-u/ddos/manifest/status"
)

func NewDaemonManifest() *daemon.Manifest {
	return new(daemon.Manifest)
}

func NewStatusManifest() *status.Manifest {
	return new(status.Manifest)
}

func NewDDoSManifest() *ddos.Manifest {
	return new(ddos.Manifest)
}
