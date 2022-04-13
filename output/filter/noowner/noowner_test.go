package noowner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"yanic/data"
	"yanic/runtime"
)

func TestFilter(t *testing.T) {
	assert := assert.New(t)

	// invalid config
	_, err := build("nope")
	assert.Error(err)

	// delete owner by configuration
	filter, err := build(true)
	assert.NoError(err)
	n := filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{
		Owner: &data.Owner{
			Contact: "blub",
		},
	}})

	assert.NotNil(n)
	assert.Nil(n.Nodeinfo.Owner)

	// keep owner configuration
	filter, _ = build(false)
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{
		Owner: &data.Owner{
			Contact: "blub",
		},
	}})

	assert.NotNil(n)
	assert.NotNil(n.Nodeinfo.Owner)
}
