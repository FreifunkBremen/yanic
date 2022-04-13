package domainassite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"yanic/data"
	"yanic/runtime"
)

func TestFilter(t *testing.T) {
	assert := assert.New(t)

	// invalid config
	_, err := build("nope")
	assert.Error(err)

	// delete owner by configuration
	filter, err := build(true)
	assert.NoError(err)
	n := filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{
		System: data.System{
			SiteCode:   "ffhb",
			DomainCode: "city",
		},
	}})

	assert.NotNil(n)
	assert.Equal("city", n.Nodeinfo.System.SiteCode)
	assert.Equal("", n.Nodeinfo.System.DomainCode)

	// keep owner configuration
	filter, _ = build(false)
	n = filter.Apply(&runtime.Node{Nodeinfo: &data.Nodeinfo{
		System: data.System{
			SiteCode:   "ffhb",
			DomainCode: "city",
		},
	}})

	assert.NotNil(n)
	assert.Equal("ffhb", n.Nodeinfo.System.SiteCode)
	assert.Equal("city", n.Nodeinfo.System.DomainCode)
}
