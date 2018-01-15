package respondd

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	assert := assert.New(t)

	conn, err := Connect(map[string]interface{}{
		"type":    "udp6",
		"address": "fasfs",
	})
	assert.Nil(conn)
	assert.Error(err)

	conn, err = Connect(map[string]interface{}{
		"type":    "udp",
		"address": "localhost:11001",
	})
	assert.NoError(err)

	conn.InsertNode(&runtime.Node{
		Nodeinfo: &data.NodeInfo{
			NodeID:   "73deadbeaf13",
			Hostname: "inject-test",
			Network: data.Network{
				Mac:       "73:de:ad:be:af:13",
				Addresses: []string{"a", "b"},
			},
		},
		Statistics: &data.Statistics{
			NodeID: "73deadbeaf13",
			Clients: data.Clients{
				Total:  1000,
				Wifi:   500,
				Wifi24: 100,
				Wifi5:  300,
			},
		},
	})

	conn.Close()

}
