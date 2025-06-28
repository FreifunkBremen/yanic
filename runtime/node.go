package runtime

import (
	"net"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/lib/jsontime"
)

// Node struct
type Node struct {
	Address      *net.UDPAddr           `json:"-"` // the last known address
	Firstseen    jsontime.Time          `json:"firstseen"`
	Lastseen     jsontime.Time          `json:"lastseen"`
	Online       bool                   `json:"online"`
	Statistics   *data.Statistics       `json:"statistics"`
	Nodeinfo     *data.Nodeinfo         `json:"nodeinfo"`
	Neighbours   *data.Neighbours       `json:"-"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

// IsGateway returns whether the node is a gateway
func (node *Node) IsGateway() bool {
	if info := node.Nodeinfo; info != nil {
		return info.VPN
	}
	return false
}
