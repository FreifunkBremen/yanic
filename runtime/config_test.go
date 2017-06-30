package runtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	assert := assert.New(t)

	config, err := ReadConfigFile("../config_example.toml")
	assert.NoError(err, "no error during reading")
	assert.NotNil(config)

	assert.True(config.Respondd.Enable)
	assert.Equal("eth0", config.Respondd.Interface)
	assert.Equal(time.Minute, config.Respondd.CollectInterval.Duration)

	assert.Equal(time.Hour*24*7, config.Nodes.PruneAfter.Duration)

	assert.Equal(time.Hour*24*7, config.Database.DeleteAfter.Duration)

	var meshviewer map[string]interface{}
	outputs := config.Nodes.Output["meshviewer"]
	assert.Len(outputs, 1, "more outputs are given")
	meshviewer = outputs[0].(map[string]interface{})
	assert.Equal(int64(2), meshviewer["version"])
	assert.Equal("/var/www/html/meshviewer/data/nodes.json", meshviewer["nodes_path"])

	var influxdb map[string]interface{}
	dbs := config.Database.Connection["influxdb"]
	assert.Len(dbs, 1, "more influxdb are given")
	influxdb = dbs[0].(map[string]interface{})
	assert.Equal(influxdb["database"], "ffhb")

	var graphitedb map[string]interface{}
	dbs = config.Database.Connection["graphite"]
	assert.Len(dbs, 1, "more graphitedb are given")
	graphitedb = dbs[0].(map[string]interface{})
	assert.Equal(graphitedb["address"], "localhost:2003")
}
