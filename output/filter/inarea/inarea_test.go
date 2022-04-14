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

	n := filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{
		Location: &data.Location{Latitude: 4.0, Longitude: 11.0},
	}})
	assert.NotNil(n)

	// keep without nodeinfo -> use has_location for it
	n = filter.Apply(&runtime.Node{})
	assert.NotNil(n)

	// keep without location -> use has_location for it
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{}})
	assert.NotNil(n)

	// zeros not in area (0, 0)
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{
		Location: &data.Location{},
	}})
	assert.Nil(n)

	// in area
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{
		Location: &data.Location{Latitude: 4.0, Longitude: 11.0},
	}})
	assert.NotNil(n)

	// over max longitude -> dropped
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{
		Location: &data.Location{Latitude: 4.0, Longitude: 13.0},
	}})
	assert.Nil(n)

	// over max latitude -> dropped
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{
		Location: &data.Location{Latitude: 6.0, Longitude: 11.0},
	}})
	assert.Nil(n)

	// lower then mix latitde -> dropped
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{
		Location: &data.Location{Latitude: 1.0, Longitude: 2.0},
	}})
	assert.Nil(n)

	// invalid config format
	_, err := build(true)
	assert.Error(err)

	// invalid config latitude switched max and min
	_, err = build(map[string]interface{}{
		"latitude_min":  5.0,
		"latitude_max":  3.0,
		"longitude_min": 10.0,
		"longitude_max": 12.0,
	})
	assert.Error(err)

	// invalid config longitude switched max and min
	_, err = build(map[string]interface{}{
		"latitude_min":  3.0,
		"latitude_max":  5.0,
		"longitude_min": 15.0,
		"longitude_max": 10.0,
	})
	assert.Error(err)

}
