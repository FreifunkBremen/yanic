package domainassite

import (
	"errors"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/output/filter"
	"github.com/FreifunkBremen/yanic/runtime"
)

type domainAsSite struct{ set bool }

func init() {
	filter.Register("domain_as_site", build)
}

func build(config interface{}) (filter.Filter, error) {
	if value, ok := config.(bool); ok {
		return &domainAsSite{set: value}, nil
	}
	return nil, errors.New("invalid configuration, boolean expected")
}

func (config *domainAsSite) Apply(node *runtime.Node) *runtime.Node {
	if nodeinfo := node.Nodeinfo; nodeinfo != nil && config.set && nodeinfo.System.DomainCode != "" {
		node = &runtime.Node{
			Address:    node.Address,
			Firstseen:  node.Firstseen,
			Lastseen:   node.Lastseen,
			Online:     node.Online,
			Statistics: node.Statistics,
			Nodeinfo: &data.Nodeinfo{
				NodeID:  nodeinfo.NodeID,
				Network: nodeinfo.Network,
				System: data.System{
					SiteCode: nodeinfo.System.DomainCode,
				},
				Owner:    nodeinfo.Owner,
				Hostname: nodeinfo.Hostname,
				Location: nodeinfo.Location,
				Software: nodeinfo.Software,
				Hardware: nodeinfo.Hardware,
				VPN:      nodeinfo.VPN,
				Wireless: nodeinfo.Wireless,
			},
			Neighbours: node.Neighbours,
		}
	}
	return node
}
