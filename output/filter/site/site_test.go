package site

import (
	"testing"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestFilterSite(t *testing.T) {
	assert := assert.New(t)

	filter, err := build([]string{"ffhb"})
	assert.NoError(err)

	n := filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{System: data.System{SiteCode: "ffxx"}}})
	assert.Nil(n)

	n = filter.Apply(&runtime.Node{Nodeinfo: &data.NodeInfo{System: data.System{SiteCode: "ffhb"}}})
	assert.NotNil(n)

	n = filter.Apply(&runtime.Node{})
	assert.Nil(n)
}
