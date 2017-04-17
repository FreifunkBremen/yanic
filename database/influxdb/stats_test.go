package influxdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
)

func TestGlobalStats(t *testing.T) {
	stats := runtime.NewGlobalStats(createTestNodes(), "ffhb01")

	assert := assert.New(t)
	fields := GlobalStatsFields(stats)

	// check fields
	assert.EqualValues(2, fields["nodes"])
}

func createTestNodes() *runtime.Nodes {
	nodes := runtime.NewNodes(&runtime.Config{})

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
			System: data.System{
				SiteCode: "ffhb01",
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
			System: data.System{
				SiteCode: "ffhb01",
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
