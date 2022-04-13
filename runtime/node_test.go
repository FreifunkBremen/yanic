package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"yanic/data"
)

func TestNode(t *testing.T) {
	assert := assert.New(t)

	node := &Node{}
	assert.False(node.IsGateway())

	node.Nodeinfo = &data.Nodeinfo{}
	assert.False(node.IsGateway())

	node.Nodeinfo.VPN = true
	assert.True(node.IsGateway())

	node.Nodeinfo.VPN = false
	assert.False(node.IsGateway())
}
