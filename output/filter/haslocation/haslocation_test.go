package haslocation

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/output/filter"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFilterHasLocation(t *testing.T) {
	assert := assert.New(t)
	var config filter.Config

	config = map[string]interface{}{}

	filterHasLocation := config.HasLocation()
	n := filterHasLocation(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Location: &data.Location{},
	}})
	assert.NotNil(n)

	config["has_location"] = true
	filterHasLocation = config.HasLocation()

	n = filterHasLocation(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Location: &data.Location{},
	}})
	assert.NotNil(n)

	n = filterHasLocation(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.Nil(n)

	n = filterHasLocation(&runtime.Node{})
	assert.Nil(n)

	config["has_location"] = false
	filterHasLocation = config.HasLocation()

	n = filterHasLocation(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Location: &data.Location{},
	}})
	assert.Nil(n)

	n = filterHasLocation(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.NotNil(n)

	n = filterHasLocation(&runtime.Node{})
	assert.NotNil(n)
}
