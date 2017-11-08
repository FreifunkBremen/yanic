package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
)

func TestGlobalStats(t *testing.T) {
	stats := NewGlobalStats(createTestNodes())

	assert := assert.New(t)
	assert.EqualValues(1, stats.Gateways)
	assert.EqualValues(3, stats.Nodes)
	assert.EqualValues(25, stats.Clients)

	// check models
	assert.Len(stats.Models, 2)
	assert.EqualValues(2, stats.Models["TP-Link 841"])
	assert.EqualValues(1, stats.Models["Xeon Multi-Core"])

	// check firmwares
	assert.Len(stats.Firmwares, 1)
	assert.EqualValues(1, stats.Firmwares["2016.1.6+entenhausen1"])

	// check autoupdater
	assert.Len(stats.Autoupdater, 2)
	assert.EqualValues(1, stats.Autoupdater["stable"])
}

func createTestNodes() *Nodes {
	nodes := NewNodes(&Config{})

	nodeData := &Node{
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

	nodes.AddNode(&Node{
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
			Software: data.Software{
				Autoupdater: struct {
					Enabled bool   `json:"enabled,omitempty"`
					Branch  string `json:"branch,omitempty"`
				}{
					Enabled: true,
					Branch:  "stable",
				},
			},
		},
	})

	nodes.AddNode(&Node{
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
