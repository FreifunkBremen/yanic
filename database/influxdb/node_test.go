package influxdb

import (
	"testing"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
)

func TestToInflux(t *testing.T) {
	assert := assert.New(t)

	node := &runtime.Node{
		Statistics: &data.Statistics{
			NodeID:      "deadbeef",
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
						"BAFF1E5": data.BatmanLink{
							Tq: 75,
						},
					},
				},
			},
			LLDP: map[string]data.LLDPNeighbours{},
		},
	}

	points := testPoints(node)
	var fields map[string]interface{}
	var tags map[string]string

	assert.Len(points, 2)

	// first point contains the neighbour
	nPoint := points[0]
	tags = nPoint.Tags()
	fields, _ = nPoint.Fields()
	assert.EqualValues("link", nPoint.Name())
	assert.EqualValues(map[string]string{"source": "deadbeef", "target": "BAFF1E5"}, tags)
	assert.EqualValues(75, fields["tq"])

	// second point contains the statistics
	sPoint := points[1]
	tags = sPoint.Tags()
	fields, _ = sPoint.Fields()

	assert.EqualValues("deadbeef", tags["nodeid"])
	assert.EqualValues("nobody", tags["owner"])
	assert.EqualValues(0.5, fields["load"])
	assert.EqualValues(0, fields["neighbours.lldp"])
	assert.EqualValues(1, fields["neighbours.batadv"])
	assert.EqualValues(1, fields["neighbours.vpn"])
	assert.EqualValues(1, fields["neighbours.total"])

	assert.EqualValues(uint32(3), fields["wireless.txpower24"])
	assert.EqualValues(uint32(5500), fields["airtime11a.frequency"])
	assert.EqualValues("", tags["frequency5500"])

	assert.EqualValues(int64(1213), fields["traffic.rx.bytes"])
	assert.EqualValues(float64(1321), fields["traffic.tx.dropped"])
	assert.EqualValues(int64(1322), fields["traffic.forward.bytes"])
	assert.EqualValues(int64(2331), fields["traffic.mgmt_rx.bytes"])
	assert.EqualValues(float64(2327), fields["traffic.mgmt_tx.packets"])
}

// Processes data and returns the InfluxDB points
func testPoints(nodes ...*runtime.Node) (points []*client.Point) {
	// Create dummy client
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{Addr: "http://127.0.0.1"})
	if err != nil {
		panic(err)
	}

	// Create dummy connection
	conn := &Connection{
		points: make(chan *client.Point),
		client: influxClient,
	}

	// Process data
	go func() {
		for _, node := range nodes {
			conn.InsertNode(node)
		}
		conn.Close()
	}()

	// Read points
	for point := range conn.points {
		points = append(points, point)
	}

	return
}
