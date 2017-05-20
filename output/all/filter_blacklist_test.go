package all

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFilterBlacklist(t *testing.T) {
	assert := assert.New(t)
	var config filterConfig

	config = map[string]interface{}{}

	filterBlacklist := config.Blacklist()

	n := filterBlacklist(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.NotNil(n)

	config["blacklist"] = []interface{}{"a", "c"}
	filterBlacklist = config.Blacklist()

	n = filterBlacklist(&runtime.Node{Nodeinfo: &data.NodeInfo{NodeID: "a"}})
	assert.Nil(n)

	n = filterBlacklist(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.NotNil(n)

	n = filterBlacklist(&runtime.Node{})
	assert.NotNil(n)

}
