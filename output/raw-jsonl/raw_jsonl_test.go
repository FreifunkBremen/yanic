package jsonlines

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
)

func TestTransform(t *testing.T) {
	testNodes := createTestNodes()
	result := transform(testNodes)

	assert := assert.New(t)
	assert.Len(testNodes.List, 4)
	assert.Len(result, 5)

	fi, ok := result[0].(FileInfo)
	assert.True(ok)
	assert.Equal(1, fi.Version)
	assert.Equal("raw-nodes-jsonl", fi.Format)

	foundNodeIDs := make(map[string]int)

	for _, element := range result[1:] {
		node, ok := element.(*RawNode)
		assert.True(ok)
		foundNodeIDs[node.Nodeinfo.NodeID] += 1
	}

	assert.Equal(1, foundNodeIDs["abcdef012425"])
	assert.Equal(1, foundNodeIDs["abcdef012345"])
	assert.Equal(1, foundNodeIDs["112233445566"])
	assert.Equal(1, foundNodeIDs["0xdeadbeef0x"])
	assert.Equal(0, foundNodeIDs["NONEXISTING"])
}

func createTestNodes() *runtime.Nodes {
	nodes := runtime.NewNodes(&runtime.NodesConfig{})

	nodes.AddNode(&runtime.Node{
		Online: true,
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 42,
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
	nodeData.Nodeinfo.Software.Firmware = &struct {
		Base      string `json:"base,omitempty"`
		Release   string `json:"release,omitempty"`
		Target    string `json:"target,omitempty"`
		Subtarget string `json:"subtarget,omitempty"`
		ImageName string `json:"image_name,omitempty"`
	}{
		Release: "2019.1~exp42",
	}
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
