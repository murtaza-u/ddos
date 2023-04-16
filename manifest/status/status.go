package status

import "github.com/murtaza-u/ddos/manifest/daemon"

type Manifest struct {
	Daemons []daemon.Manifest `json:"daemons,omitempty"`
}
