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
}

func createTestNodes() *Nodes {
	nodes := NewNodes(&Config{})

	nodeData := &data.ResponseData{
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 23,
			},
		},
		NodeInfo: &data.NodeInfo{
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
		},
	}
	nodeData.NodeInfo.Software.Firmware.Release = "2016.1.6+entenhausen1"
	nodes.Update("abcdef012345", nodeData)

	nodes.Update("112233445566", &data.ResponseData{
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 2,
			},
		},
		NodeInfo: &data.NodeInfo{
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
		},
	})

	nodes.Update("0xdeadbeef0x", &data.ResponseData{
		NodeInfo: &data.NodeInfo{
			VPN: true,
			Hardware: data.Hardware{
				Model: "Xeon Multi-Core",
			},
		},
	})

	return nodes
}
