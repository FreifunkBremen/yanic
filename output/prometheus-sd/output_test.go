package prometheus_sd

import (
	"os"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"yanic/data"
	"yanic/runtime"
)

func TestOutput(t *testing.T) {
	assert := assert.New(t)

	out, err := Register(map[string]interface{}{})
	assert.Error(err)
	assert.Nil(out)

	nodes := runtime.NewNodes(&runtime.NodesConfig{})
	ipAddress, err := net.ResolveUDPAddr("udp6", "[fe80::20de:a:3ac%eth0]:1001")
	assert.NoError(err)
	nodes.AddNode(&runtime.Node{
		Online:  true,
		Address: ipAddress,
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

	// IP
	out, err = Register(map[string]interface{}{
		"path": "/tmp/prometheus_sd.json",
	})
	os.Remove("/tmp/prometheus_sd.json")
	assert.NoError(err)
	assert.NotNil(out)

	out.Save(nodes)
	_, err = os.Stat("/tmp/prometheus_sd.json")
	assert.NoError(err)

	// NodeID
	out, err = Register(map[string]interface{}{
		"target_address": "node_id",
		"path":           "/tmp/prometheus_sd.json",
		"labels": map[string]interface{}{
			"hosts":   "ffhb",
			"service": "yanic",
		},
	})
	os.Remove("/tmp/prometheus_sd.json")
	assert.NoError(err)
	assert.NotNil(out)

	out.Save(nodes)
	_, err = os.Stat("/tmp/prometheus_sd.json")
	assert.NoError(err)
}
