package meshviewerFFRGB

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	assert := assert.New(t)
	nodes := runtime.NewNodes(&runtime.Config{})
	node := NewNode(nodes, &runtime.Node{
		Nodeinfo: &data.NodeInfo{
			Network: data.Network{
				Mac: "blub",
			},
			Location: &data.Location{
				Longtitude: 13.3,
				Latitude:   8.7,
			},
		},
		Statistics: &data.Statistics{
			Memory: data.Memory{
				Free:  13,
				Total: 50,
			},
			Wireless: []*data.WirelessAirtime{
				&data.WirelessAirtime{
					ChanUtil:  0.3,
					Frequency: 2512,
				},
				&data.WirelessAirtime{
					ChanUtil:  0.4,
					Frequency: 2612,
				},
				&data.WirelessAirtime{
					ChanUtil:  0.5,
					Frequency: 5200,
				},
			},
		},
	})
	assert.NotNil(node)
	assert.Equal("blub", node.Network.MAC)
	assert.Equal(13.3, node.Location.Longtitude)
	assert.Equal(8.7, node.Location.Latitude)
	assert.Equal(0.74, *node.MemoryUsage)
}
