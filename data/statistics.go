package data

/*
	Nodes Lua based respondd do not have a integer type.
	They always return float.
*/

type Statistics struct {
	NodeId      string  `json:"node_id"`
	Clients     Clients `json:"clients"`
	RootFsUsage float64 `json:"rootfs_usage,omitempty"`
	LoadAverage float64 `json:"loadavg,omitempty"`
	Memory      Memory  `json:"memory,omitempty"`
	Uptime      float64 `json:"uptime,omitempty"`
	Idletime    float64 `json:"idletime,omitempty"`
	Gateway     string  `json:"gateway,omitempty"`
	Processes   struct {
		Total   uint32 `json:"total"`
		Running uint32 `json:"running"`
	} `json:"processes,omitempty"`
	MeshVpn *MeshVPN `json:"mesh_vpn,omitempty"`
	Traffic struct {
		Tx      *Traffic `json:"tx"`
		Rx      *Traffic `json:"rx"`
		Forward *Traffic `json:"forward"`
		MgmtTx  *Traffic `json:"mgmt_tx"`
		MgmtRx  *Traffic `json:"mgmt_rx"`
	} `json:"traffic,omitempty"`
	Switch   map[string]*SwitchPort `json:"switch,omitempty"`
	Wireless *WirelessStatistics    `json:"wireless,omitempty"`
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
	Cached  uint32 `json:"cached"`
	Total   uint32 `json:"total"`
	Buffers uint32 `json:"buffers"`
	Free    uint32 `json:"free"`
}

type SwitchPort struct {
	Speed uint32 `json:"speed"`
}
