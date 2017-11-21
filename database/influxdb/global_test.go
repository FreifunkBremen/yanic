package influxdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
)

const TEST_SITE = "ffxx"

func TestGlobalStats(t *testing.T) {
	stats := runtime.NewGlobalStats(createTestNodes(), []string{TEST_SITE})

	assert := assert.New(t)

	// check SITE_GLOBAL fields
	fields := GlobalStatsFields(stats[runtime.GLOBAL_SITE])
	assert.EqualValues(3, fields["nodes"])

	// check TEST_SITE fields
	fields = GlobalStatsFields(stats[TEST_SITE])
	assert.EqualValues(2, fields["nodes"])
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
			System: data.System{
				SiteCode: TEST_SITE,
			},
		},
	}
	nodeData.Nodeinfo.Software.Firmware.Release = "2016.1.6+entenhausen1"
	nodeData.Nodeinfo.Software.Autoupdater.Enabled = true
	nodeData.Nodeinfo.Software.Autoupdater.Branch = "stable"
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
			System: data.System{
				SiteCode: TEST_SITE,
			},
		},
	})

	return nodes
}
