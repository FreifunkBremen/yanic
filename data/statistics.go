package data

/*
	Nodes Lua based respondd do not have a integer type.
	They always return float.
*/

type Statistics struct {
	NodeId      string   `json:"node_id"`
	Clients     Clients  `json:"clients"`
	RootFsUsage float64  `json:"rootfs_usage"`
	Traffic     *Traffic `json:"traffic"`
	Memory      Memory   `json:"memory"`
	Uptime      float64  `json:"uptime"`
	Idletime    float64  `json:"idletime"`
	Gateway     string   `json:"gateway"`
	Processes   struct {
		Total   uint32 `json:"total"`
		Running uint32 `json:"running"`
	} `json:"processes"`
	LoadAverage float64  `json:"loadavg"`
	MeshVpn     *MeshVPN `json:"mesh_vpn,omitempty"`
}

type MeshVPNPeerLink struct {
	Established float64 `json:"established"`
}

type MeshVPNPeerGroup struct {
	Peers  map[string]*MeshVPNPeerLink  `json:"peers"`
	Groups map[string]*MeshVPNPeerGroup `json:"groups"`
}

type MeshVPN struct {
	Groups map[string]*MeshVPNPeerGroup `json:"groups,omitempty"`
}

type TrafficEntry struct {
	Bytes   float64 `json:"bytes,omitempty"`
	Packets float64 `json:"packets,omitempty"`
	Dropped float64 `json:"dropped,omitempty"`
}

type Traffic struct {
	Tx      *TrafficEntry `json:"tx"`
	Rx      *TrafficEntry `json:"rx"`
	Forward *TrafficEntry `json:"forward"`
	MgmtTx  *TrafficEntry `json:"mgmt_tx"`
	MgmtRx  *TrafficEntry `json:"mgmt_rx"`
}

type Clients struct {
	Wifi   uint32 `json:"wifi"`
	Wifi24 uint32 `json:"wifi24"`
	Wifi5  uint32 `json:"wifi5"`
	Total  uint32 `json:"total"`
}

type Memory struct {
	Cached  uint64 `json:"cached"`
	Total   uint64 `json:"total"`
	Buffers uint64 `json:"buffers"`
	Free    uint64 `json:"free"`
}
