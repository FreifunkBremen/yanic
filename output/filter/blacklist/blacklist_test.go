package blacklist

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFilterBlacklist(t *testing.T) {
	assert := assert.New(t)

	// invalid config
	filter, err := build(3)
	assert.Error(err)

	// tests with empty list
	filter, err = build([]string{})

	// keep node without nodeid
	n := filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.NotNil(n)

	// tests with blacklist
	filter, _ = build([]string{"a", "c"})

	// blacklist contains node with nodeid -> drop it
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{NodeID: "a"}})
	assert.Nil(n)

	// blacklist does not contains node without nodeid -> keep it
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.NotNil(n)

	n = filter.Apply(&runtime.Node{})
	assert.NotNil(n)
}
