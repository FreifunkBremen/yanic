package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	assert := assert.New(t)

	config, err := ReadConfigFile("../config_example.toml")
	assert.NoError(err)
	assert.NotNil(config)

	assert.True(config.Respondd.Enable)
	assert.Equal("br-ffhb", config.Respondd.Interfaces[0].InterfaceName)
	assert.Equal(time.Minute, config.Respondd.CollectInterval.Duration)
	assert.Equal(time.Hour*24*7, config.Nodes.PruneAfter.Duration)
	assert.Equal(time.Hour*24*7, config.Database.DeleteAfter.Duration)

	assert.Len(config.Respondd.Sites, 1)
	assert.Contains(config.Respondd.Sites, "ffhb")
	assert.Contains(config.Respondd.Sites["ffhb"].Domains, "city")

	// Test output plugins
	assert.Len(config.Nodes.Output, 5)
	outputs := config.Nodes.Output["meshviewer"].([]interface{})
	assert.Len(outputs, 1)
	meshviewer := outputs[0]

	assert.EqualValues(map[string]interface{}{
		"version":    int64(2),
		"enable":     false,
		"nodes_path": "/var/www/html/meshviewer/data/nodes.json",
		"graph_path": "/var/www/html/meshviewer/data/graph.json",
		"filter": map[string]interface{}{
			"no_owner": true,
		},
	}, meshviewer)

	_, err = ReadConfigFile("testdata/config_invalid.toml")
	assert.Error(err, "not unmarshalable")
	assert.Contains(err.Error(), "invalid TOML syntax")

	_, err = ReadConfigFile("testdata/adsa.toml")
	assert.Error(err, "not found able")
	assert.Contains(err.Error(), "no such file or directory")
}
