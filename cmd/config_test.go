package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	config, err := ReadConfigFile("../config_example.toml")
	require.NoError(err)
	require.NotNil(config)

	assert.True(config.Respondd.Enable)
	assert.Equal("br-ffhb", config.Respondd.Interfaces[0].InterfaceName)
	assert.Equal(time.Minute, config.Respondd.CollectInterval.Duration)
	assert.Equal(time.Hour*24*7, config.Nodes.PruneAfter.Duration)
	assert.Equal(time.Hour*24*7, config.Database.DeleteAfter.Duration)

	assert.Len(config.Respondd.Sites, 1)
	assert.Contains(config.Respondd.Sites, "ffhb")
	assert.Contains(config.Respondd.Sites["ffhb"].Domains, "city")

	// Test output plugins
	assert.Len(config.Nodes.Output, 6)
	outputs := config.Nodes.Output["meshviewer"].([]map[string]interface{})
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
	assert.EqualError(err, "toml: line 2: expected '.' or '=', but got '\\n' instead")

	_, err = ReadConfigFile("testdata/adsa.toml")
	assert.EqualError(err, "open testdata/adsa.toml: no such file or directory")
}
