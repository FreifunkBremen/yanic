package models

import (
	"testing"

	"github.com/FreifunkBremen/respond-collector/data"
	"github.com/stretchr/testify/assert"
)

func TestToInflux(t *testing.T) {
	assert := assert.New(t)

	node := Node{
		Statistics: &data.Statistics{
			NodeId:      "foobar",
			LoadAverage: 0.5,
		},
		Nodeinfo: &data.NodeInfo{
			Owner: &data.Owner{
				Contact: "nobody",
			},
		},
		Neighbours: &data.Neighbours{},
	}

	tags, fields := node.ToInflux()

	assert.Equal("foobar", tags.GetString("nodeid"))
	assert.Equal("nobody", tags.GetString("owner"))
	assert.Equal(0.5, fields["load"])
}
