package haslocation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"yanic/data"
	"yanic/runtime"
)

func TestFilterHasLocation(t *testing.T) {
	assert := assert.New(t)

	// invalid config
	_, err := build(3)
	assert.Error(err)

	// test to drop nodes without location
	filter, err := build(true)
	assert.NoError(err)

	// node has location (with 0,0) -> keep it
	n := filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{
		Location: &data.Location{},
	}})
	assert.NotNil(n)

	// node without location has no location -> drop it
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{}})
	assert.Nil(n)

	// node without nodeinfo has no location -> drop it
	n = filter.Apply(&runtime.Node{})
	assert.Nil(n)

	// test to drop nodes without location
	filter, err = build(false)
	assert.NoError(err)

	// node has location (with 0,0) -> drop it
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{
		Location: &data.Location{},
	}})
	assert.Nil(n)

	// node without location has no location -> keep it
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{}})
	assert.NotNil(n)

	// node without nodeinfo has no location -> keep it
	n = filter.Apply(&runtime.Node{})
	assert.NotNil(n)
}
