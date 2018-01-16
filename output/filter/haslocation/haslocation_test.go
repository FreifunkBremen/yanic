package haslocation

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFilterHasLocation(t *testing.T) {
	assert := assert.New(t)

	filter, err := build(true)
	assert.NoError(err)

	n := filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Location: &data.Location{},
	}})
	assert.NotNil(n)

	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.Nil(n)

	n = filter.Apply(&runtime.Node{})
	assert.Nil(n)

	filter, err = build(false)
	assert.NoError(err)

	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Location: &data.Location{},
	}})
	assert.Nil(n)

	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.NotNil(n)

	n = filter.Apply(&runtime.Node{})
	assert.NotNil(n)
}
