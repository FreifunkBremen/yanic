package geojson

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
)

const (
	testNodeDescription string = "Online; 23 Clients; Model: TP-Link 841;" +
		" Site: mysite; Domain: domain_42;"
)

func TestTransform(t *testing.T) {
	testNodes := createTestNodes()
	nodes := transform(testNodes)

	assert := assert.New(t)
	assert.Len(testNodes.List, 4)
	assert.Len(nodes.Features, 3)

	node := testNodes.List["abcdef012425"]

	umap := getUMapOptions(node)
	assert.Len(umap, 2)

	description := getNodeDescription(node)
	assert.Len(description, len(testNodeDescription))
}

func createTestNodes() *runtime.Nodes {
	nodes := runtime.NewNodes(&runtime.NodesConfig{})

	nodes.AddNode(&runtime.Node{
		Online: true,
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 23,
			},
		},
		Nodeinfo: &data.Nodeinfo{
			NodeID: "abcdef012425",
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
			Location: &data.Location{
				Latitude:  24,
				Longitude: 2,
			},
			System: data.System{
				SiteCode:   "mysite",
				DomainCode: "domain_42",
			},
		},
	})

	nodeData := &runtime.Node{
		Online: true,
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 23,
			},
		},
		Nodeinfo: &data.Nodeinfo{
			NodeID: "abcdef012345",
			Hardware: data.Hardware{
				Model: "TP-Link 842",
			},
			System: data.System{
				SiteCode:   "mysite",
				DomainCode: "domain_42",
			},
		},
	}
	nodeData.Nodeinfo.Software.Firmware.Release = "2019.1~exp42"
	nodes.AddNode(nodeData)

	nodes.AddNode(&runtime.Node{
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 2,
			},
		},
		Nodeinfo: &data.Nodeinfo{
			NodeID: "112233445566",
			Hardware: data.Hardware{
				Model: "TP-Link 843",
			},
			Location: &data.Location{
				Latitude:  23,
				Longitude: 2,
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
			Location: &data.Location{
				Latitude:  23,
				Longitude: 22,
			},
		},
	})

	return nodes
}
