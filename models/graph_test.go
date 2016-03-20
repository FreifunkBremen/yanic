package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateGraph(t *testing.T) {
	assert := assert.New(t)

	nodes := &Nodes{
		List: make(map[string]*Node),
	}

	graph := nodes.BuildGraph()
	assert.NotNil(graph)
	// TODO more tests required
}
