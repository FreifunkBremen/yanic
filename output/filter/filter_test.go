package filter

import (
	"errors"
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

type filterBool struct{ bool }

func (f filterBool) Apply(node *runtime.Node) *runtime.Node {
	if f.bool {
		return node
	}
	return nil
}

func build(v interface{}) (Filter, error) {
	if config, ok := v.(bool); ok {
		return &filterBool{config}, nil
	}
	return nil, nil
}

func buildNil(v interface{}) (Filter, error) {
	return nil, nil
}

func buildError(v interface{}) (Filter, error) {
	if v != nil {
		return nil, errors.New("some errors")
	}
	return nil, nil
}

func TestFilter(t *testing.T) {
	assert := assert.New(t)

	Register("test_nil", buildNil)
	Register("test_err", buildError)
	Register("test", build)

	// filter still exists
	filter, err := New(map[string]interface{}{
		"adsa": true,
	})
	assert.Len(err, 1)
	assert.Nil(filter)

	// no filter
	filter, err = New(map[string]interface{}{
		"test_nil": 3,
	})
	assert.Len(err, 0)
	assert.Nil(filter)

	// filter error
	filter, err = New(map[string]interface{}{
		"test_err": false,
	})
	assert.Len(err, 1)
	assert.Nil(filter)

	// filter a node
	nodes := &runtime.Nodes{
		List: map[string]*runtime.Node{
			"a": {
				Nodeinfo: &data.NodeInfo{NodeID: "a"},
			},
		},
	}
	filter, err = New(map[string]interface{}{
		"test": false,
	})
	assert.Len(err, 0)
	nodes = filter.Apply(nodes)
	assert.Len(nodes.List, 0)

	// keep a node
	nodes = &runtime.Nodes{
		List: map[string]*runtime.Node{
			"a": {
				Nodeinfo: &data.NodeInfo{NodeID: "a"},
			},
		},
	}
	filter, err = New(map[string]interface{}{
		"test": true,
	})
	assert.Len(err, 0)
	nodes = filter.Apply(nodes)
	assert.Len(nodes.List, 1)
}
