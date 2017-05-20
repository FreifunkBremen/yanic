package all

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFilterNoOwner(t *testing.T) {
	assert := assert.New(t)
	var config filterConfig

	config = map[string]interface{}{}

	filterNoOwner := config.NoOwner()
	n := filterNoOwner(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Owner: &data.Owner{
			Contact: "blub",
		},
	}})
	assert.NotNil(n)
	assert.Nil(n.Nodeinfo.Owner)

	n = filterNoOwner(&runtime.Node{})
	assert.NotNil(n)

	config["no_owner"] = true
	filterNoOwner = config.NoOwner()
	n = filterNoOwner(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Owner: &data.Owner{
			Contact: "blub",
		},
	}})
	assert.NotNil(n)
	assert.Nil(n.Nodeinfo.Owner)

	config["no_owner"] = false
	filterNoOwner = config.NoOwner()

	n = filterNoOwner(&runtime.Node{Nodeinfo: &data.NodeInfo{
		Owner: &data.Owner{
			Contact: "blub",
		},
	}})
	assert.NotNil(n)
	assert.NotNil(n.Nodeinfo.Owner)
}
