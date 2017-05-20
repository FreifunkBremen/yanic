package influxdb

import (
	"testing"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/influxdata/influxdb/models"

	"github.com/stretchr/testify/assert"
)

func TestAddPoint(t *testing.T) {
	assert := assert.New(t)

	// Test add Point with tags
	connection := &Connection{
		config: map[string]interface{}{
			"tags": map[string]interface{}{
				"testtag": "value",
			},
		},
		points: make(chan *client.Point, 1),
	}

	connection.addPoint("name", models.Tags{}, models.Fields{"clients.total": 10}, time.Now())
	point := <-connection.points
	assert.NotNil(point)
	tags := point.Tags()
	assert.NotNil(tags)
	assert.Equal(tags["testtag"], "value")
	assert.NotEqual(tags["testtag2"], "value")

	// Tried to overright by config
	connection.config["tags"] = map[string]interface{}{
		"nodeid": "value",
	}

	tagsOrigin := models.Tags{}
	tagsOrigin.SetString("nodeid", "collected")

	connection.addPoint("name", tagsOrigin, models.Fields{"clients.total": 10}, time.Now())
	point = <-connection.points
	assert.NotNil(point)
	tags = point.Tags()
	assert.NotNil(tags)
	assert.Equal(tags["nodeid"], "collected")

	// Test panic if it was not possible to create a point
	assert.Panics(func() {
		connection.addPoint("name", models.Tags{}, nil, time.Now())
	})
}
