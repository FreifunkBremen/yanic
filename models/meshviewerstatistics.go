package models

import "github.com/FreifunkBremen/respond-collector/data"

type MeshviewerStatistics struct {
	NodeId      string  `json:"node_id"`
	Clients     uint32 `json:"clients"`
	RootFsUsage float64 `json:"rootfs_usage,omitempty""`
	LoadAverage float64 `json:"loadavg,omitempty""`
	Memory      data.Memory  `json:"memory,omitempty""`
	Uptime      float64 `json:"uptime,omitempty""`
	Idletime    float64 `json:"idletime,omitempty""`
	Gateway     string  `json:"gateway,omitempty"`
	Processes   struct {
		Total   uint32 `json:"total"`
		Running uint32 `json:"running"`
	} `json:"processes,omitempty""`
	MeshVpn *data.MeshVPN `json:"mesh_vpn,omitempty"`
	Traffic struct {
		Tx      *data.Traffic `json:"tx"`
		Rx      *data.Traffic `json:"rx"`
		Forward *data.Traffic `json:"forward"`
		MgmtTx  *data.Traffic `json:"mgmt_tx"`
		MgmtRx  *data.Traffic `json:"mgmt_rx"`
	} `json:"traffic,omitempty""`
}
