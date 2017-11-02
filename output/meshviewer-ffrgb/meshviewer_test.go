package meshviewerFFRGB

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	assert := assert.New(t)

	nodes := runtime.NewNodes(&runtime.Config{})
	nodes.AddNode(&runtime.Node{
		Online: true,
		Nodeinfo: &data.NodeInfo{
			NodeID: "node_a",
			Network: data.Network{
				Mac: "node:a:mac",
				Mesh: map[string]*data.BatInterface{
					"bat0": &data.BatInterface{
						Interfaces: struct {
							Wireless []string `json:"wireless,omitempty"`
							Other    []string `json:"other,omitempty"`
							Tunnel   []string `json:"tunnel,omitempty"`
						}{
							Wireless: []string{"node:a:mac:wifi"},
							Tunnel:   []string{"node:a:mac:vpn"},
							Other:    []string{"node:a:mac:lan"},
						},
					},
				},
			},
		},
		Neighbours: &data.Neighbours{
			NodeID: "node_a",
			Batadv: map[string]data.BatadvNeighbours{
				"node:a:mac:wifi": data.BatadvNeighbours{
					Neighbours: map[string]data.BatmanLink{
						"node:b:mac:wifi": data.BatmanLink{Tq: 153},
					},
				},
				"node:a:mac:lan": data.BatadvNeighbours{
					Neighbours: map[string]data.BatmanLink{
						"node:b:mac:lan": data.BatmanLink{Tq: 51},
					},
				},
			},
		},
	})
	nodes.AddNode(&runtime.Node{
		Online: true,
		Nodeinfo: &data.NodeInfo{
			NodeID: "node_c",
			Network: data.Network{
				Mac: "node:c:mac",
				Mesh: map[string]*data.BatInterface{
					"bat0": &data.BatInterface{
						Interfaces: struct {
							Wireless []string `json:"wireless,omitempty"`
							Other    []string `json:"other,omitempty"`
							Tunnel   []string `json:"tunnel,omitempty"`
						}{
							Other: []string{"node:c:mac:lan"},
						},
					},
				},
			},
		},
		Neighbours: &data.Neighbours{
			NodeID: "node_b",
			Batadv: map[string]data.BatadvNeighbours{
				"node:c:mac:lan": data.BatadvNeighbours{
					Neighbours: map[string]data.BatmanLink{
						"node:b:mac:lan": data.BatmanLink{Tq: 102},
					},
				},
			},
		},
	})
	nodes.AddNode(&runtime.Node{
		Online: true,
		Nodeinfo: &data.NodeInfo{
			NodeID: "node_b",
			Network: data.Network{
				Mac: "node:b:mac",
				Mesh: map[string]*data.BatInterface{
					"bat0": &data.BatInterface{
						Interfaces: struct {
							Wireless []string `json:"wireless,omitempty"`
							Other    []string `json:"other,omitempty"`
							Tunnel   []string `json:"tunnel,omitempty"`
						}{
							Wireless: []string{"node:b:mac:wifi"},
							Other:    []string{"node:b:mac:lan"},
						},
					},
				},
			},
		},
		Neighbours: &data.Neighbours{
			NodeID: "node_b",
			Batadv: map[string]data.BatadvNeighbours{
				"node:b:mac:lan": data.BatadvNeighbours{
					Neighbours: map[string]data.BatmanLink{
						"node:c:mac:lan": data.BatmanLink{Tq: 204},
					},
				},
				"node:b:mac:wifi": data.BatadvNeighbours{
					Neighbours: map[string]data.BatmanLink{
						"node:a:mac:wifi": data.BatmanLink{Tq: 204},
					},
				},
			},
		},
	})
	nodes.AddNode(&runtime.Node{
		Online: false,
		Nodeinfo: &data.NodeInfo{
			NodeID: "node_d",
			Network: data.Network{
				Mac: "node:d:mac",
				Mesh: map[string]*data.BatInterface{
					"bat0": &data.BatInterface{
						Interfaces: struct {
							Wireless []string `json:"wireless,omitempty"`
							Other    []string `json:"other,omitempty"`
							Tunnel   []string `json:"tunnel,omitempty"`
						}{
							Wireless: []string{"node:b:mac:wifi"},
							Other:    []string{"node:b:mac:lan"},
						},
					},
				},
			},
		},
		Neighbours: &data.Neighbours{
			NodeID: "node_d",
			Batadv: map[string]data.BatadvNeighbours{
				"node:d:mac:lan": data.BatadvNeighbours{
					Neighbours: map[string]data.BatmanLink{
						"node:c:mac:lan": data.BatmanLink{Tq: 204},
					},
				},
				"node:d:mac:wifi": data.BatadvNeighbours{
					Neighbours: map[string]data.BatmanLink{
						"node:a:mac:wifi": data.BatmanLink{Tq: 204},
					},
				},
			},
		},
	})

	meshviewer := transform(nodes)
	assert.NotNil(meshviewer)
	assert.Len(meshviewer.Nodes, 4)
	links := meshviewer.Links
	assert.Len(links, 3)

	for _, link := range links {
		switch link.SourceMAC {
		case "node:a:mac:lan":
			assert.Equal("other", link.Type)
			assert.Equal("node:b:mac:lan", link.TargetMAC)
			assert.Equal(float32(0.2), link.SourceTQ)
			assert.Equal(float32(0.2), link.TargetTQ)
			break

		case "node:a:mac:wifi":
			assert.Equal("wifi", link.Type)
			assert.Equal("node:b:mac:wifi", link.TargetMAC)
			assert.Equal(float32(0.6), link.SourceTQ)
			assert.Equal(float32(0.8), link.TargetTQ)
		default:
			assert.Equal("other", link.Type)
			assert.Equal("node:c:mac:lan", link.TargetMAC)
			assert.Equal(float32(0.8), link.SourceTQ)
			assert.Equal(float32(0.4), link.TargetTQ)
			break
		}
	}
}