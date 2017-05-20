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
	assert.Equal([]string{"eth0"}, config.Respondd.Interfaces)
	assert.Equal(time.Minute, config.Respondd.CollectInterval.Duration)

	assert.Equal(time.Hour*24*7, config.Nodes.PruneAfter.Duration)

	assert.Equal(time.Hour*24*7, config.Database.DeleteAfter.Duration)

	var meshviewer map[string]interface{}
	var outputs []map[string]interface{}
	outputs = config.Nodes.Output["meshviewer"].([]map[string]interface{})
	assert.Len(outputs, 1, "more outputs are given")
	meshviewer = outputs[0]
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

	_, err = ReadConfigFile("testdata/config_failed.toml")
	assert.Error(err, "not unmarshalable")
	assert.Contains(err.Error(), "Near line ")

	_, err = ReadConfigFile("testdata/adsa.toml")
	assert.Error(err, "not found able")
	assert.Contains(err.Error(), "no such file or directory")
}
