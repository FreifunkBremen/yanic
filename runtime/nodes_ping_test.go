package runtime

import (
	"net"
	"testing"
	"time"

	"github.com/bdlm/log"
	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
)

func TestPing(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	assert := assert.New(t)
	config := &NodesConfig{
		PingCount: 1,
	}
	config.OfflineAfter.Duration = time.Minute * 10
	// to get default (100%) path of testing
	// config.PruneAfter.Duration = time.Hour * 24 * 6
	nodes := &Nodes{
		config:        config,
		List:          make(map[string]*Node),
		ifaceToNodeID: make(map[string]string),
	}

	node := nodes.Update("expire", &data.ResponseData{NodeInfo: &data.NodeInfo{
		NodeID:  "nodeID-Lola",
		Network: data.Network{Addresses: []string{"fe80::1", "fd2f::1"}},
	}})
	// get fallback
	assert.False(nodes.ping(node))

	node.Address = &net.UDPAddr{Zone: "bat0"}
	// error during ping
	assert.False(nodes.ping(node))

	node.Address.IP = net.ParseIP("fe80::1")
	// error during ping
	assert.False(nodes.ping(node))
}
