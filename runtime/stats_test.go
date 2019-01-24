package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
)

const (
	TEST_SITE   = "ffhb"
	TEST_DOMAIN = "city"
)

func TestGlobalStats(t *testing.T) {
	stats := NewGlobalStats(createTestNodes(), map[string][]string{TEST_SITE: {TEST_DOMAIN}})

	assert := assert.New(t)
	assert.Len(stats, 2)

	//check GLOBAL_SITE stats
	assert.EqualValues(1, stats[GLOBAL_SITE][GLOBAL_DOMAIN].Gateways)
	assert.EqualValues(3, stats[GLOBAL_SITE][GLOBAL_DOMAIN].Nodes)
	assert.EqualValues(25, stats[GLOBAL_SITE][GLOBAL_DOMAIN].Clients)

	// check models
	assert.Len(stats[GLOBAL_SITE][GLOBAL_DOMAIN].Models, 2)
	assert.EqualValues(2, stats[GLOBAL_SITE][GLOBAL_DOMAIN].Models["TP-Link 841"])
	assert.EqualValues(1, stats[GLOBAL_SITE][GLOBAL_DOMAIN].Models["Xeon Multi-Core"])

	// check firmwares
	assert.Len(stats[GLOBAL_SITE][GLOBAL_DOMAIN].Firmwares, 1)
	assert.EqualValues(1, stats[GLOBAL_SITE][GLOBAL_DOMAIN].Firmwares["2016.1.6+entenhausen1"])

	// check autoupdater
	assert.Len(stats[GLOBAL_SITE][GLOBAL_DOMAIN].Autoupdater, 2)
	assert.EqualValues(1, stats[GLOBAL_SITE][GLOBAL_DOMAIN].Autoupdater["stable"])

	// check TEST_SITE stats
	assert.EqualValues(1, stats[TEST_SITE][GLOBAL_DOMAIN].Gateways)
	assert.EqualValues(2, stats[TEST_SITE][GLOBAL_DOMAIN].Nodes)
	assert.EqualValues(23, stats[TEST_SITE][GLOBAL_DOMAIN].Clients)

	// check models
	assert.Len(stats[TEST_SITE][GLOBAL_DOMAIN].Models, 2)
	assert.EqualValues(1, stats[TEST_SITE][GLOBAL_DOMAIN].Models["TP-Link 841"])
	assert.EqualValues(1, stats[TEST_SITE][GLOBAL_DOMAIN].Models["Xeon Multi-Core"])

	// check firmwares
	assert.Len(stats[TEST_SITE][GLOBAL_DOMAIN].Firmwares, 1)
	assert.EqualValues(1, stats[TEST_SITE][GLOBAL_DOMAIN].Firmwares["2016.1.6+entenhausen1"])

	// check autoupdater
	assert.Len(stats[TEST_SITE][GLOBAL_DOMAIN].Autoupdater, 1)
	assert.EqualValues(0, stats[TEST_SITE][GLOBAL_DOMAIN].Autoupdater["stable"])

	// check TEST_DOMAIN stats
	assert.EqualValues(1, stats[TEST_SITE][TEST_DOMAIN].Gateways)
	assert.EqualValues(1, stats[TEST_SITE][TEST_DOMAIN].Nodes)
	assert.EqualValues(0, stats[TEST_SITE][TEST_DOMAIN].Clients)

	// check models
	assert.Len(stats[TEST_SITE][TEST_DOMAIN].Models, 1)
	assert.EqualValues(0, stats[TEST_SITE][TEST_DOMAIN].Models["TP-Link 841"])
	assert.EqualValues(1, stats[TEST_SITE][TEST_DOMAIN].Models["Xeon Multi-Core"])

	// check firmwares
	assert.Len(stats[TEST_SITE][TEST_DOMAIN].Firmwares, 0)
	assert.EqualValues(0, stats[TEST_SITE][TEST_DOMAIN].Firmwares["2016.1.6+entenhausen1"])

	// check autoupdater
	assert.Len(stats[TEST_SITE][TEST_DOMAIN].Autoupdater, 1)
	assert.EqualValues(0, stats[TEST_SITE][TEST_DOMAIN].Autoupdater["stable"])
}

func createTestNodes() *Nodes {
	nodes := NewNodes(&NodesConfig{})

	nodeData := &Node{
		Online: true,
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
			System: data.System{
				SiteCode: TEST_SITE,
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
		Nodeinfo: &data.Nodeinfo{
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
		Nodeinfo: &data.Nodeinfo{
			NodeID: "0xdeadbeef0x",
			VPN:    true,
			Hardware: data.Hardware{
				Model: "Xeon Multi-Core",
			},
			System: data.System{
				SiteCode:   TEST_SITE,
				DomainCode: TEST_DOMAIN,
			},
		},
	})

	return nodes
}
