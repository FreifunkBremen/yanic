package meshviewerFFRGB

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	assert := assert.New(t)

	nodes := runtime.NewNodes(&runtime.NodesConfig{})
	nodes.AddNode(&runtime.Node{
		Online: true,
		Nodeinfo: &data.Nodeinfo{
			NodeID: "node_a",
			Network: data.Network{
				Mac: "node:a:mac",
				Mesh: map[string]*data.NetworkInterface{
					"bat0": {
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
				"node:a:mac:wifi": {
					Neighbours: map[string]data.BatmanLink{
						"node:b:mac:wifi": {Tq: 153},
					},
				},
				"node:a:mac:lan": {
					Neighbours: map[string]data.BatmanLink{
						"node:b:mac:lan": {Tq: 51},
					},
				},
			},
		},
	})
	nodes.AddNode(&runtime.Node{
		Online: true,
		Nodeinfo: &data.Nodeinfo{
			NodeID: "node_c",
			Network: data.Network{
				Mac: "node:c:mac",
				Mesh: map[string]*data.NetworkInterface{
					"bat0": {
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
			NodeID: "node_c",
			Batadv: map[string]data.BatadvNeighbours{
				"node:c:mac:lan": {
					Neighbours: map[string]data.BatmanLink{
						"node:b:mac:lan": {Tq: 102},
					},
				},
			},
		},
	})
	nodes.AddNode(&runtime.Node{
		Online: true,
		Nodeinfo: &data.Nodeinfo{
			NodeID: "node_b",
			Network: data.Network{
				Mac: "node:b:mac",
				Mesh: map[string]*data.NetworkInterface{
					"bat0": {
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
				"node:b:mac:lan": {
					Neighbours: map[string]data.BatmanLink{
						"node:c:mac:lan": {Tq: 204},
					},
				},
				"node:b:mac:wifi": {
					Neighbours: map[string]data.BatmanLink{
						"node:a:mac:wifi": {Tq: 204},
					},
				},
			},
		},
	})
	nodes.AddNode(&runtime.Node{
		Online: false,
		Nodeinfo: &data.Nodeinfo{
			NodeID: "node_d",
			Network: data.Network{
				Mac: "node:d:mac",
				Mesh: map[string]*data.NetworkInterface{
					"bat0": {
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
				"node:d:mac:lan": {
					Neighbours: map[string]data.BatmanLink{
						"node:c:mac:lan": {Tq: 204},
					},
				},
				"node:d:mac:wifi": {
					Neighbours: map[string]data.BatmanLink{
						"node:a:mac:wifi": {Tq: 204},
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
		switch link.SourceAddress {
		case "node:a:mac:lan":
			assert.Equal("other", link.Type)
			assert.Equal("node:b:mac:lan", link.TargetAddress)
			assert.Equal(float32(0.2), link.SourceTQ)
			assert.Equal(float32(0), link.TargetTQ)
			break

		case "node:a:mac:wifi":
			assert.Equal("wifi", link.Type)
			assert.Equal("node:b:mac:wifi", link.TargetAddress)
			assert.Equal(float32(0.6), link.SourceTQ)
			assert.Equal(float32(0.8), link.TargetTQ)
		case "node:b:mac:lan":
			assert.Equal("other", link.Type)
			assert.Equal("node:c:mac:lan", link.TargetAddress)
			assert.Equal(float32(0.8), link.SourceTQ)
			assert.Equal(float32(0.4), link.TargetTQ)
			break
		default:
			assert.False(true, "invalid link.SourceAddress found")
		}
	}
}
