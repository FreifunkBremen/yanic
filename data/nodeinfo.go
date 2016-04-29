package data

type NodeInfo struct {
	NodeId   string    `json:"node_id"`
	Network  Network   `json:"network"`
	Owner    *Owner    `json:"owner,omitempty"`
	System   System    `json:"system"`
	Hostname string    `json:"hostname"`
	Location *Location `json:"location,omitempty"`
	Software Software  `json:"software"`
	Hardware Hardware  `json:"hardware"`
	VPN      bool      `json:"vpn"`
}
type BatInterface struct {
	Interfaces struct {
		Wireless []string `json:"wireless,omitempty"`
		Other    []string `json:"other,omitempty"`
		Tunnel   []string `json:"tunnel,omitempty"`
	} `json:"interfaces"`
}

type Network struct {
	Mac       string   `json:"mac"`
	Addresses []string `json:"addresses"`
	Mesh      map[string]*BatInterface `json:"mesh"`
	MeshInterfaces []string `json:"mesh_interfaces"`
}

type Owner struct {
	Contact string `json:"contact"`
}

type System struct {
	SiteCode string `json:"site_code"`
}

type Location struct {
	Longtitude float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	Altitude   float64 `json:"altitude,omitempty"`
}

type Software struct {
	Autoupdater *struct {
		Enabled bool   `json:"enabled"`
		Branch  string `json:"branch"`
	} `json:"autoupdater,omitempty"`
	BatmanAdv *struct {
		Version string `json:"version"`
		Compat  int    `json:"compat"`
	} `json:"batman-adv,omitempty"`
	Fastd *struct {
		Enabled bool   `json:"enabled"`
		Version string `json:"version"`
	} `json:"fastd,omitempty"`
	Firmware *struct {
		Base    string `json:"base"`
		Release string `json:"release"`
	} `json:"firmware,omitempty"`
	StatusPage *struct {
		Api int `json:"api"`
	} `json:"status-page,omitempty"`
}

type Hardware struct {
	Nproc int    `json:"nproc"`
	Model string `json:"model"`
}
