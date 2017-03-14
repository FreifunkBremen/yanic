package data

/*
	Nodes Lua based respondd do not have a integer type.
	They always return float.
*/

//Statistics struct
type Statistics struct {
	NodeID         string   `json:"node_id"`
	Clients        *Clients `json:"clients"`
	RootFsUsage    float64  `json:"rootfs_usage,omitempty"`
	LoadAverage    float64  `json:"loadavg,omitempty"`
	Memory         *Memory  `json:"memory,omitempty"`
	Uptime         float64  `json:"uptime,omitempty"`
	Idletime       float64  `json:"idletime,omitempty"`
	GatewayIPv4    string   `json:"gateway,omitempty"`
	GatewayIPv6    string   `json:"gateway6,omitempty"`
	GatewayNexthop string   `json:"gateway_nexthop,omitempty"`
	Processes      *struct {
		Total   uint32 `json:"total"`
		Running uint32 `json:"running"`
	} `json:"processes,omitempty"`
	MeshVPN *MeshVPN `json:"mesh_vpn,omitempty"`
	Traffic struct {
		Tx      *Traffic `json:"tx"`
		Rx      *Traffic `json:"rx"`
		Forward *Traffic `json:"forward"`
		MgmtTx  *Traffic `json:"mgmt_tx"`
		MgmtRx  *Traffic `json:"mgmt_rx"`
	} `json:"traffic,omitempty"`
	Switch   map[string]*SwitchPort `json:"switch,omitempty"`
	Wireless WirelessStatistics     `json:"wireless,omitempty"`
}

// MeshVPNPeerLink struct
type MeshVPNPeerLink struct {
	Established float64 `json:"established"`
}

// MeshVPNPeerGroup struct
type MeshVPNPeerGroup struct {
	Peers  map[string]*MeshVPNPeerLink  `json:"peers"`
	Groups map[string]*MeshVPNPeerGroup `json:"groups"`
}

// MeshVPN struct
type MeshVPN struct {
	Groups map[string]*MeshVPNPeerGroup `json:"groups,omitempty"`
}

// Traffic struct
type Traffic struct {
	Bytes   float64 `json:"bytes,omitempty"`
	Packets float64 `json:"packets,omitempty"`
	Dropped float64 `json:"dropped,omitempty"`
}

// Clients struct
type Clients struct {
	Wifi   uint32 `json:"wifi"`
	Wifi24 uint32 `json:"wifi24"`
	Wifi5  uint32 `json:"wifi5"`
	Total  uint32 `json:"total"`
}

// Memory struct
type Memory struct {
	Cached  uint32 `json:"cached"`
	Total   uint32 `json:"total"`
	Buffers uint32 `json:"buffers"`
	Free    uint32 `json:"free"`
}

// SwitchPort struct
type SwitchPort struct {
	Speed uint32 `json:"speed"`
}
