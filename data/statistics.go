package data

/*
	Nodes Lua based respondd do not have a integer type.
	They always return float.
*/

//Statistics struct
type Statistics struct {
	NodeID         string  `json:"node_id"`
	Clients        Clients `json:"clients"`
	DHCP           *DHCP   `json:"dhcp"`
	RootFsUsage    float64 `json:"rootfs_usage,omitempty"`
	LoadAverage    float64 `json:"loadavg,omitempty"`
	Memory         Memory  `json:"memory,omitempty"`
	Uptime         float64 `json:"uptime,omitempty"`
	Idletime       float64 `json:"idletime,omitempty"`
	GatewayIPv4    string  `json:"gateway,omitempty"`
	GatewayIPv6    string  `json:"gateway6,omitempty"`
	GatewayNexthop string  `json:"gateway_nexthop,omitempty"`
	Processes      struct {
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
	Switch    map[string]*SwitchPort `json:"switch,omitempty"`
	Wireless  WirelessStatistics     `json:"wireless,omitempty"`
	ProcStats *ProcStats             `json:"stat,omitempty"`
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
	Owe    uint32 `json:"owe"`
	Owe24  uint32 `json:"owe24"`
	Owe5   uint32 `json:"owe5"`
	Total  uint32 `json:"total"`
}

// DHCP struct
type DHCP struct {
	// Packet counters
	Decline  uint32 `json:"dhcp_decline"`
	Offer    uint32 `json:"dhcp_offer"`
	Ack      uint32 `json:"dhcp_ack"`
	Nak      uint32 `json:"dhcp_nak"`
	Request  uint32 `json:"dhcp_request"`
	Discover uint32 `json:"dhcp_discover"`
	Inform   uint32 `json:"dhcp_inform"`
	Release  uint32 `json:"dhcp_release"`

	LeasesAllocated uint32 `json:"leases_allocated_4"`
	LeasesPruned    uint32 `json:"leases_pruned_4"`
}

// Memory struct
type Memory struct {
	Cached    int64 `json:"cached"`
	Total     int64 `json:"total"`
	Buffers   int64 `json:"buffers"`
	Free      int64 `json:"free,omitempty"`
	Available int64 `json:"available,omitempty"`
}

// SwitchPort struct
type SwitchPort struct {
	Speed uint32 `json:"speed"`
}

// ProcStats struct
type ProcStats struct {
	CPU             ProcStatsCPU `json:"cpu"`
	Intr            int64        `json:"intr"`
	ContextSwitches int64        `json:"ctxt"`
	SoftIRQ         int64        `json:"softirq"`
	Processes       int64        `json:"processes"`
}

// ProcStatsCPU struct
type ProcStatsCPU struct {
	User    int64 `json:"user"`
	Nice    int64 `json:"nice"`
	System  int64 `json:"system"`
	Idle    int64 `json:"idle"`
	IOWait  int64 `json:"iowait"`
	IRQ     int64 `json:"irq"`
	SoftIRQ int64 `json:"softirq"`
}
