package runtime

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/stretchr/testify/assert"
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
