package all

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFilterSites(t *testing.T) {
	assert := assert.New(t)
	var config filterConfig

	config = map[string]interface{}{}

	filterSites := config.Sites()

	n := filterSites(&runtime.Node{Nodeinfo: &data.NodeInfo{}})
	assert.NotNil(n)

	config["sites"] = []interface{}{"ffhb"}
	filterSites = config.Sites()

	n = filterSites(&runtime.Node{Nodeinfo: &data.NodeInfo{System: data.System{SiteCode: "ffxx"}}})
	assert.Nil(n)

	n = filterSites(&runtime.Node{Nodeinfo: &data.NodeInfo{System: data.System{SiteCode: "ffhb"}}})
	assert.NotNil(n)

	n = filterSites(&runtime.Node{})
	assert.Nil(n)

}
