package meshviewer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
)

func TestNodesV1(t *testing.T) {
	nodes := BuildNodesV1(createTestNodes()).(*NodesV1)

	assert := assert.New(t)
	assert.Len(nodes.List, 2)
}
func TestNodesV2(t *testing.T) {
	nodes := BuildNodesV2(createTestNodes()).(*NodesV2)

	assert := assert.New(t)
	assert.Len(nodes.List, 2)
}

func createTestNodes() *runtime.Nodes {
	nodes := runtime.NewNodes(&runtime.NodesConfig{})

	nodeData := &runtime.Node{
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 23,
			},
		},
		Nodeinfo: &data.Nodeinfo{
			NodeID: "abcdef012345",
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
		},
	}
	nodeData.Nodeinfo.Software.Firmware.Release = "2016.1.6+entenhausen1"
	nodes.AddNode(nodeData)

	nodes.AddNode(&runtime.Node{
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Wifi24: 2,
				Wifi5:  3,
			},
			Memory: data.Memory{
				Total: 32,
				Free:  8,
			},
		},
		Nodeinfo: &data.Nodeinfo{
			NodeID: "112233445566",
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
		},
	})

	nodes.AddNode(&runtime.Node{
		Nodeinfo: &data.Nodeinfo{
			NodeID: "0xdeadbeef0x",
			VPN:    true,
			Hardware: data.Hardware{
				Model: "Xeon Multi-Core",
			},
		},
	})

	return nodes
}
