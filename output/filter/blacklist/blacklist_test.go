package blacklist

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFilterBlacklist(t *testing.T) {
	assert := assert.New(t)

	filter, _ := build([]string{})

	n := filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.NotNil(n)

	filter, _ = build([]string{"a", "c"})
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{NodeID: "a"}})
	assert.Nil(n)

	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.NotNil(n)

	n = filter.Apply(&runtime.Node{})
	assert.NotNil(n)
}
