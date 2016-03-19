package data

/*
	Nodes Lua based respondd do not have a integer type.
	They always return float.
*/

type Statistics struct {
	NodeId      string  `json:"node_id"`
	Clients     Clients `json:"clients"`
	RootFsUsage float64 `json:"rootfs_usage"`
	LoadAverage float64 `json:"loadavg"`
	Memory      Memory  `json:"memory"`
	Uptime      float64 `json:"uptime"`
	Idletime    float64 `json:"idletime"`
	Gateway     string  `json:"gateway"`
	Processes   struct {
		Total   uint32 `json:"total"`
		Running uint32 `json:"running"`
	} `json:"processes"`
	MeshVpn *MeshVPN `json:"mesh_vpn,omitempty"`
	Traffic struct {
		Tx      *Traffic `json:"tx"`
		Rx      *Traffic `json:"rx"`
		Forward *Traffic `json:"forward"`
		MgmtTx  *Traffic `json:"mgmt_tx"`
		MgmtRx  *Traffic `json:"mgmt_rx"`
	} `json:"traffic"`
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

type Traffic struct {
	Bytes   float64 `json:"bytes,omitempty"`
	Packets float64 `json:"packets,omitempty"`
	Dropped float64 `json:"dropped,omitempty"`
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
