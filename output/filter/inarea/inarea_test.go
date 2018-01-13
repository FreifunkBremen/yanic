package inarea

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFilterInArea(t *testing.T) {
	assert := assert.New(t)

	filter, _ := build(map[string]interface{}{
		"latitude_min":  3.0,
		"latitude_max":  5.0,
		"longitude_min": 10.0,
		"longitude_max": 12.0,
	})

	n := filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Location: &data.Location{Latitude: 4.0, Longitude: 11.0},
	}})
	assert.NotNil(n)

	// drop area without nodeinfo
	n = filter.Apply(&runtime.Node{})
	assert.Nil(n)

	// keep without location
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.NotNil(n)

	// zeros not in area (0, 0)
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Location: &data.Location{},
	}})
	assert.Nil(n)

	// in area
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Location: &data.Location{Latitude: 4.0, Longitude: 11.0},
	}})
	assert.NotNil(n)

	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Location: &data.Location{Latitude: 4.0, Longitude: 13.0},
	}})
	assert.Nil(n)

	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Location: &data.Location{Latitude: 6.0, Longitude: 11.0},
	}})
	assert.Nil(n)

	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Location: &data.Location{Latitude: 1.0, Longitude: 2.0},
	}})
	assert.Nil(n)
}
