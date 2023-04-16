package daemon

import "time"

type Manifest struct {
	Metadata Metadata `json:"metadata,omitempty"`
}

type Metadata struct {
	Name     string    `json:"name,omitempty"`
	Alias    string    `json:"alias,omitempty"`
	Info     Info      `json:"info,omitempty"`
	LastPing time.Time `json:"lastPing,omitempty"`
}

type Info struct {
	Hostname string `json:"hostname,omitempty"`
	Username string `json:"username,omitempty"`
	OS       string `json:"os,omitempty"`
	Platform string `json:"platform,omitempty"`
	Cores    uint64 `json:"cores,omitempty"`
	Memory   uint64 `json:"memory,omitempty"`
}
