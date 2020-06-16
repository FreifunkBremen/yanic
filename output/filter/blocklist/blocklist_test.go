package blocklist

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFilterBlocklist(t *testing.T) {
	assert := assert.New(t)

	// invalid config
	filter, err := build(3)
	assert.Error(err)

	filter, err = build([]interface{}{2, "a"})
	assert.Error(err)

	// tests with empty list
	filter, err = build([]interface{}{})
	assert.NoError(err)

	// keep node without nodeid
	n := filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{}})
	assert.NotNil(n)

	// tests with blocklist
	filter, err = build([]interface{}{"a", "c"})
	assert.NoError(err)

	// blocklist contains node with nodeid -> drop it
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{NodeID: "a"}})
	assert.Nil(n)

	// blocklist does not contains node without nodeid -> keep it
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{}})
	assert.NotNil(n)

	n = filter.Apply(&runtime.Node{})
	assert.NotNil(n)
}
