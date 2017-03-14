package meshviewer

import (
	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/jsontime"
)

// Node struct
type Node struct {
	Firstseen  jsontime.Time    `json:"firstseen"`
	Lastseen   jsontime.Time    `json:"lastseen"`
	Flags      Flags            `json:"flags"`
	Statistics *Statistics      `json:"statistics"`
	Nodeinfo   *data.NodeInfo   `json:"nodeinfo"`
	Neighbours *data.Neighbours `json:"-"`
}

// Flags status of node set by collector for the meshviewer
type Flags struct {
	Online  bool `json:"online"`
	Gateway bool `json:"gateway"`
}

// Statistics a meshviewer spezifisch struct, diffrent from respondd
type Statistics struct {
	NodeID      string   `json:"node_id"`
	Clients     uint32   `json:"clients"`
	RootFsUsage float64  `json:"rootfs_usage,omitempty"`
	LoadAverage float64  `json:"loadavg,omitempty"`
	MemoryUsage *float64 `json:"memory_usage,omitempty"`
	Uptime      float64  `json:"uptime,omitempty"`
	Idletime    float64  `json:"idletime,omitempty"`
	GatewayIPv4 string   `json:"gateway,omitempty"`
	GatewayIPv6 string   `json:"gateway6,omitempty"`
	Processes   *struct {
		Total   uint32 `json:"total"`
		Running uint32 `json:"running"`
	} `json:"processes,omitempty"`
	MeshVPN *data.MeshVPN `json:"mesh_vpn,omitempty"`
	Traffic struct {
		Tx      *data.Traffic `json:"tx"`
		Rx      *data.Traffic `json:"rx"`
		Forward *data.Traffic `json:"forward"`
		MgmtTx  *data.Traffic `json:"mgmt_tx"`
		MgmtRx  *data.Traffic `json:"mgmt_rx"`
	} `json:"traffic,omitempty"`
}

// NewStatistics transform respond Statistics to meshviewer Statistics
func NewStatistics(stats *data.Statistics) *Statistics {
	var total uint32
	if clients := stats.Clients; clients != nil {
		total = clients.Total
		if total <= 0 {
			total = clients.Wifi24 + clients.Wifi5
		}
	}
	/* The Meshviewer could not handle absolute memory output
	 * calc the used memory as a float which 100% equal 1.0
	 * calc is coppied from node statuspage (look discussion:
	 * https://github.com/FreifunkBremen/yanic/issues/35)
	 */

	meshviewerStats := &Statistics{
		NodeID:      stats.NodeID,
		GatewayIPv4: stats.GatewayIPv4,
		GatewayIPv6: stats.GatewayIPv6,
		RootFsUsage: stats.RootFsUsage,
		LoadAverage: stats.LoadAverage,
		Uptime:      stats.Uptime,
		Idletime:    stats.Idletime,
		Processes:   stats.Processes,
		MeshVPN:     stats.MeshVPN,
		Traffic:     stats.Traffic,
		Clients:     total,
	}
	if memory := stats.Memory; memory != nil && memory.Total > 0 {
		*meshviewerStats.MemoryUsage = 1 - (float64(memory.Free)+float64(memory.Buffers)+float64(memory.Cached))/float64(memory.Total)
	}
	return meshviewerStats
}
