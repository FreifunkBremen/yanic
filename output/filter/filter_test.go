package filter

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	assert := assert.New(t)

	// filtered - do not run all
	nodes := &runtime.Nodes{
		List: map[string]*runtime.Node{
			"a": {
				Nodeinfo: &data.NodeInfo{NodeID: "a"},
			},
		},
	}
	config := filterConfig{
		"has_location": true,
	}
	nodes = config.filtering(nodes)
	assert.Len(nodes.List, 0)

	// run to end
	nodes = &runtime.Nodes{
		List: map[string]*runtime.Node{
			"a": {
				Nodeinfo: &data.NodeInfo{NodeID: "a"},
			},
		},
	}
	config = filterConfig{
		"has_location": false,
	}
	nodes = config.filtering(nodes)
	assert.Len(nodes.List, 1)
}
