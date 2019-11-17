package meshviewerFFRGB

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	assert := assert.New(t)
	nodes := runtime.NewNodes(&runtime.NodesConfig{})
	node := NewNode(nodes, &runtime.Node{
		Nodeinfo: &data.Nodeinfo{
			Owner: &data.Owner{
				Contact: "whoami",
			},
			Network: data.Network{
				Mac:       "blub",
				Addresses: []string{"fe80::1"},
			},
		},
	})
	assert.NotNil(node)
	assert.Len(node.Addresses, 1)

	node = NewNode(nodes, &runtime.Node{
		Nodeinfo: &data.Nodeinfo{
			Owner: &data.Owner{
				Contact: "whoami",
			},
			Network: data.Network{
				Mac: "blub",
			},
			Location: &data.Location{
				Longitude: 13.3,
				Latitude:  8.7,
			},
		},
		Statistics: &data.Statistics{
			Memory: data.Memory{
				Free:  13,
				Total: 50,
			},
			Wireless: []*data.WirelessAirtime{
				{
					ChanUtil:  0.3,
					Frequency: 2512,
				},
				{
					ChanUtil:  0.4,
					Frequency: 2612,
				},
				{
					ChanUtil:  0.5,
					Frequency: 5200,
				},
			},
		},
		CustomFields: map[string]interface{}{
			"custom_fields": "are_custom",
			"custom_int":    3,
		},
	})
	assert.NotNil(node)
	assert.NotNil(node.Addresses)
	assert.Equal("whoami", node.Owner)
	assert.Equal("blub", node.MAC)
	assert.Equal(13.3, node.Location.Longitude)
	assert.Equal(8.7, node.Location.Latitude)
	assert.Equal(0.74, *node.MemoryUsage)
	assert.Equal("are_custom", node.CustomFields["custom_fields"])
	assert.Equal(3, node.CustomFields["custom_int"])
}
