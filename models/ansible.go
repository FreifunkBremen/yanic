package models

type Ansible struct {
	Nodes []string `json:"nodes"`
	Meta  struct {
		HostVars []*AnsibleHostVars `json:"hostvars"`
	} `json:"_meta"`
}
type AnsibleHostVars struct {
	Address      string  `json:"ansible_ssh_host"`
	Hostname     string  `json:"node_name"`
	Channel24    int     `json:"radio24_channel"`
	TxPower24    int     `json:"radio24_txpower"`
	Channel5     int     `json:"radio5_channel"`
	TxPower5     int     `json:"radio5_txpower"`
	GeoLatitude  float64 `json:"geo_latitude"`
	GeoLongitude float64 `json:"geo_longitude"`
}

func GenerateAnsible(nodes *Nodes, aliases map[string]*Alias) *Ansible {
	ansible := &Ansible{Nodes: make([]string, 0)}
	for nodeid, alias := range aliases {
		if node := nodes.List[nodeid]; node != nil {

			ansible.Nodes = append(ansible.Nodes, nodeid)

			vars := &AnsibleHostVars{
				Address:      node.Nodeinfo.Network.Addresses[0],
				Hostname:     alias.Hostname,
				Channel24:    alias.Freq24.Channel,
				Channel5:     alias.Freq5.Channel,
				TxPower24:    alias.Freq24.TxPower,
				TxPower5:     alias.Freq5.TxPower,
				GeoLatitude:  alias.Location.Latitude,
				GeoLongitude: alias.Location.Longtitude,
			}
			ansible.Meta.HostVars = append(ansible.Meta.HostVars, vars)

		}
	}
	return ansible
}
