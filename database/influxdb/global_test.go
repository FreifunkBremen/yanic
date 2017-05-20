package influxdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
)

func TestGlobalStats(t *testing.T) {
	stats := runtime.NewGlobalStats(createTestNodes())

	assert := assert.New(t)
	fields := GlobalStatsFields(stats)

	// check fields
	assert.EqualValues(3, fields["nodes"])
}

func createTestNodes() *runtime.Nodes {
	nodes := runtime.NewNodes(&runtime.Config{})

	nodeData := &runtime.Node{
		Online: true,
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 23,
			},
		},
		Nodeinfo: &data.NodeInfo{
			NodeID: "abcdef012345",
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
		},
	}
	nodeData.Nodeinfo.Software.Firmware.Release = "2016.1.6+entenhausen1"
	nodes.AddNode(nodeData)

	nodes.AddNode(&runtime.Node{
		Online: true,
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 2,
			},
		},
		Nodeinfo: &data.NodeInfo{
			NodeID: "112233445566",
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
		},
	})

	nodes.AddNode(&runtime.Node{
		Online: true,
		Nodeinfo: &data.NodeInfo{
			NodeID: "0xdeadbeef0x",
			VPN:    true,
			Hardware: data.Hardware{
				Model: "Xeon Multi-Core",
			},
		},
	})

	return nodes
}
