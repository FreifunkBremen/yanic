package influxdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
)

func TestToInflux(t *testing.T) {
	assert := assert.New(t)

	node := &runtime.Node{
		Statistics: &data.Statistics{
			NodeID:      "foobar",
			LoadAverage: 0.5,
			Wireless: data.WirelessStatistics{
				&data.WirelessAirtime{Frequency: 5500},
			},
			Traffic: struct {
				Tx      *data.Traffic `json:"tx"`
				Rx      *data.Traffic `json:"rx"`
				Forward *data.Traffic `json:"forward"`
				MgmtTx  *data.Traffic `json:"mgmt_tx"`
				MgmtRx  *data.Traffic `json:"mgmt_rx"`
			}{
				Tx:      &data.Traffic{Dropped: 1321},
				Rx:      &data.Traffic{Bytes: 1213},
				Forward: &data.Traffic{Bytes: 1322},
				MgmtTx:  &data.Traffic{Packets: 2327},
				MgmtRx:  &data.Traffic{Bytes: 2331},
			},
			MeshVPN: &data.MeshVPN{
				Groups: map[string]*data.MeshVPNPeerGroup{
					"ffhb": &data.MeshVPNPeerGroup{
						Peers: map[string]*data.MeshVPNPeerLink{
							"vpn01": &data.MeshVPNPeerLink{Established: 3},
							"vpn02": &data.MeshVPNPeerLink{},
							"trash": nil,
							"vpn03": &data.MeshVPNPeerLink{Established: 0},
						},
					},
				},
			},
		},
		Nodeinfo: &data.NodeInfo{
			Owner: &data.Owner{
				Contact: "nobody",
			},
			Wireless: &data.Wireless{
				TxPower24: 3,
				Channel24: 4,
			},
		},
		Neighbours: &data.Neighbours{
			Batadv: map[string]data.BatadvNeighbours{
				"a-interface": data.BatadvNeighbours{
					Neighbours: map[string]data.BatmanLink{
						"b-neigbourinterface": data.BatmanLink{},
					},
				},
			},
			LLDP: map[string]data.LLDPNeighbours{},
		},
	}

	tags, fields := nodeToInflux(node)

	assert.Equal("foobar", tags.GetString("nodeid"))
	assert.Equal("nobody", tags.GetString("owner"))
	assert.Equal(0.5, fields["load"])
	assert.Equal(0, fields["neighbours.lldp"])
	assert.Equal(1, fields["neighbours.batadv"])
	assert.Equal(1, fields["neighbours.vpn"])
	assert.Equal(1, fields["neighbours.total"])

	assert.Equal(uint32(3), fields["wireless.txpower24"])
	assert.Equal(uint32(5500), fields["airtime11a.frequency"])
	assert.Equal("", tags.GetString("frequency5500"))

	assert.Equal(int64(1213), fields["traffic.rx.bytes"])
	assert.Equal(float64(1321), fields["traffic.tx.dropped"])
	assert.Equal(int64(1322), fields["traffic.forward.bytes"])
	assert.Equal(int64(2331), fields["traffic.mgmt_rx.bytes"])
	assert.Equal(float64(2327), fields["traffic.mgmt_tx.packets"])
}
