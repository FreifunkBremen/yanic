package respond

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

const (
	SITE_TEST   = "ffhb"
	DOMAIN_TEST = "city"
)

func TestCollector(t *testing.T) {
	nodes := runtime.NewNodes(&runtime.NodesConfig{})
	config := &Config{
		Sites: map[string]SiteConfig{
			SITE_TEST: {
				Domains: []string{DOMAIN_TEST},
			},
		},
	}

	collector := NewCollector(nil, nodes, config)
	collector.Start(time.Millisecond)
	time.Sleep(time.Millisecond * 10)
	collector.Close()
}

func TestParse(t *testing.T) {
	assert := assert.New(t)

	// read testdata
	compressed, err := ioutil.ReadFile("testdata/nodeinfo.flated")
	assert.Nil(err)

	res := &Response{
		Raw: compressed,
	}

	data, err := res.parse()

	assert.NoError(err)
	assert.NotNil(data)

	assert.Equal("f81a67a5e9c1", data.Nodeinfo.NodeID)
}
