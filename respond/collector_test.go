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

	data, err := res.parse([]CustomFieldConfig{})

	assert.NoError(err)
	assert.NotNil(data)

	assert.Equal("f81a67a5e9c1", data.Nodeinfo.NodeID)
}

func TestParseCustomFields(t *testing.T) {
	assert := assert.New(t)

	// read testdata
	compressed, err := ioutil.ReadFile("testdata/nodeinfo.flated")
	assert.Nil(err)

	res := &Response{
		Raw: compressed,
	}

	customFields := []CustomFieldConfig{
		{
			Name: "my_custom_field",
			Path: "nodeinfo.hostname",
		},
	}

	data, err := res.parse(customFields)

	assert.NoError(err)
	assert.NotNil(data)

	assert.Equal("Trillian", data.CustomFields["my_custom_field"])
	assert.Equal("Trillian", data.Nodeinfo.Hostname)
}

func TestParseCustomFieldNotExistant(t *testing.T) {
	assert := assert.New(t)

	// read testdata
	compressed, err := ioutil.ReadFile("testdata/nodeinfo.flated")
	assert.Nil(err)

	res := &Response{
		Raw: compressed,
	}

	customFields := []CustomFieldConfig{
		{
			Name: "some_other_field",
			Path: "nodeinfo.some_field_which_doesnt_exist",
		},
	}

	data, err := res.parse(customFields)

	assert.NoError(err)
	assert.NotNil(data)

	_, ok := data.CustomFields["some_other_field"]
	assert.Equal("Trillian", data.Nodeinfo.Hostname)
	assert.False(ok)
}
