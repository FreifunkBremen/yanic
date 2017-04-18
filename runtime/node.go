package runtime

import (
	"net"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/jsontime"
)

// Node struct
type Node struct {
	Address    net.IP           `json:"address"` // the last known IP address
	Firstseen  jsontime.Time    `json:"firstseen"`
	Lastseen   jsontime.Time    `json:"lastseen"`
	Online     bool             `json:"online"`
	Statistics *data.Statistics `json:"statistics"`
	Nodeinfo   *data.NodeInfo   `json:"nodeinfo"`
	Neighbours *data.Neighbours `json:"-"`
}

// IsGateway returns whether the node is a gateway
func (node *Node) IsGateway() bool {
	if info := node.Nodeinfo; info != nil {
		return info.VPN
	}
	return false
}
