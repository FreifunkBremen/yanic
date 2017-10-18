package runtime

import (
	"net"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/jsontime"
)

// Node struct
type Node struct {
	Address    *net.UDPAddr     `json:"-"` // the last known address
	Firstseen  jsontime.Time    `json:"firstseen"`
	Lastseen   jsontime.Time    `json:"lastseen"`
	Online     bool             `json:"online"`
	Statistics *data.Statistics `json:"statistics"`
	Nodeinfo   *data.NodeInfo   `json:"nodeinfo"`
	Neighbours *data.Neighbours `json:"-"`
}

// Link represents a link between two nodes
type Link struct {
	SourceID  string
	SourceMAC string
	TargetID  string
	TargetMAC string
	TQ        int
}

// IsGateway returns whether the node is a gateway
func (node *Node) IsGateway() bool {
	if info := node.Nodeinfo; info != nil {
		return info.VPN
	}
	return false
}
