package models

type Ansible struct {
  Nodes []string `json:"nodes"`
  Meta struct {
    HostVars []*AnsibleHostVars `json:"hostvars"`
  } `json:"_meta"`
}
type AnsibleHostVars struct {
  Address string `json:"ansible_ssh_host"`
  Hostname string `json:"node_name"`
}

func GenerateAnsible(nodes *Nodes,aliases map[string]*Alias) *Ansible{
  ansible := &Ansible{Nodes:make([]string,0)}
  for nodeid,alias := range aliases{
    if node := nodes.List[nodeid]; node != nil {

      ansible.Nodes = append(ansible.Nodes,nodeid)

      vars := &AnsibleHostVars{
        Address: node.Nodeinfo.Network.Addresses[0],
        Hostname: alias.Hostname,
      }
      ansible.Meta.HostVars = append(ansible.Meta.HostVars,vars)

    }
  }
  return ansible
}
