package webserver

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
)

func TestPrometheusMetricsFromNode(t *testing.T) {
	assert := assert.New(t)

	m := PrometheusMetricsFromNode(nil, &runtime.Node{})
	assert.Len(m, 0)

	nodes := runtime.NewNodes(&runtime.NodesConfig{})
	nodes.AddNode(&runtime.Node{
		Nodeinfo: &data.Nodeinfo{
			NodeID: "lola",
			Network: data.Network{
				Mesh: map[string]*data.NetworkInterface{
					"mesh1": {
						Interfaces: struct {
							Wireless []string `json:"wireless,omitempty"`
							Other    []string `json:"other,omitempty"`
							Tunnel   []string `json:"tunnel,omitempty"`
						}{
							Tunnel: []string{"fe80::2"},
						},
					},
				},
			},
		},
	})

	node := &runtime.Node{
		Online: false,
		Nodeinfo: &data.Nodeinfo{
			NodeID: "wasd1",
			Network: data.Network{
				Mesh: map[string]*data.NetworkInterface{
					"mesh0": {
						Interfaces: struct {
							Wireless []string `json:"wireless,omitempty"`
							Other    []string `json:"other,omitempty"`
							Tunnel   []string `json:"tunnel,omitempty"`
						}{
							Tunnel: []string{"fe80::1"},
						},
					},
				},
			},
			Software: data.Software{
				Autoupdater: struct {
					Enabled bool   `json:"enabled,omitempty"`
					Branch  string `json:"branch,omitempty"`
				}{
					Enabled: true,
					Branch:  "testing",
				},
			},
		},
		Statistics: &data.Statistics{},
		Neighbours: &data.Neighbours{
			NodeID: "wasd1",
			Babel: map[string]data.BabelNeighbours{
				"mesh0": {
					LinkLocalAddress: "fe80::1",
					Neighbours: map[string]data.BabelLink{
						"fe80::2": {Cost: 20000},
					},
				},
			},
		},
	}
	nodes.AddNode(node)
	m = PrometheusMetricsFromNode(nodes, node)
	assert.Len(m, 15)
	assert.Equal(m[0].Labels["source_id"], "wasd1")

	m = PrometheusMetricsFromNode(nil, &runtime.Node{
		Online: true,
		Nodeinfo: &data.Nodeinfo{
			NodeID: "wasd",
			System: data.System{
				SiteCode:   "ffhb",
				DomainCode: "city",
			},
			Owner: &data.Owner{Contact: "mailto:blub@example.org"},
			Location: &data.Location{
				Latitude:  52.0,
				Longitude: 4.0,
			},
			Wireless: &data.Wireless{
				TxPower24: 0,
			},
		},
		Statistics: &data.Statistics{
			ProcStats: &data.ProcStats{},
			Traffic: struct {
				Tx      *data.Traffic `json:"tx"`
				Rx      *data.Traffic `json:"rx"`
				Forward *data.Traffic `json:"forward"`
				MgmtTx  *data.Traffic `json:"mgmt_tx"`
				MgmtRx  *data.Traffic `json:"mgmt_rx"`
			}{
				Tx:      &data.Traffic{},
				Rx:      &data.Traffic{},
				Forward: &data.Traffic{},
				MgmtTx:  &data.Traffic{},
				MgmtRx:  &data.Traffic{},
			},
			Wireless: data.WirelessStatistics{
				&data.WirelessAirtime{Frequency: 5002},
				&data.WirelessAirtime{Frequency: 2430},
			},
		},
	})

	assert.Len(m, 48)
	assert.Equal(m[0].Labels["node_id"], "wasd")
}
