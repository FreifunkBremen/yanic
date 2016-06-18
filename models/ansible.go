package models

type Ansible struct {
	Nodes []string `json:"nodes"`
	Meta  struct {
		HostVars map[string]*AnsibleHostVars `json:"hostvars"`
	} `json:"_meta"`
}
type AnsibleHostVars struct {
	Address      string  `json:"ansible_ssh_host"`
	Hostname     string  `json:"node_name,omitempty"`
	Channel24    uint32  `json:"radio24_channel,omitempty"`
	TxPower24    uint32  `json:"radio24_txpower,omitempty"`
	Channel5     uint32  `json:"radio5_channel,omitempty"`
	TxPower5     uint32  `json:"radio5_txpower,omitempty"`
	GeoLatitude  float64 `json:"geo_latitude,omitempty"`
	GeoLongitude float64 `json:"geo_longitude,omitempty"`
}

func GenerateAnsible(nodes *Nodes, aliases map[string]*Alias) *Ansible {
	ansible := &Ansible{Nodes: make([]string, 0)}
	ansible.Meta.HostVars = make(map[string]*AnsibleHostVars)
	for nodeid, alias := range aliases {
		if node := nodes.List[nodeid]; node != nil {

			ansible.Nodes = append(ansible.Nodes, nodeid)

			vars := &AnsibleHostVars{
				Address:  node.Nodeinfo.Network.Addresses[0],
				Hostname: alias.Hostname,
			}
			if alias.Freq24 != nil {
				vars.Channel24 = alias.Freq24.Channel
				vars.TxPower24 = alias.Freq24.TxPower
			}
			if alias.Freq5 != nil {
				vars.Channel5 = alias.Freq5.Channel
				vars.TxPower5 = alias.Freq5.TxPower
			}
			if alias.Location != nil {
				vars.GeoLatitude = alias.Location.Latitude
				vars.GeoLongitude = alias.Location.Longtitude
			}
			ansible.Meta.HostVars[nodeid] = vars

		}
	}
	return ansible
}
