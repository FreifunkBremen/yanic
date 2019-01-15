package influxdb

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/influxdata/influxdb1-client/models"
	"github.com/influxdata/influxdb1-client/v2"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	assert := assert.New(t)

	conn, err := Connect(map[string]interface{}{
		"address":              "",
		"username":             "",
		"password":             "",
		"insecure_skip_verify": true,
	})
	assert.Nil(conn)
	assert.Error(err)

	conn, err = Connect(map[string]interface{}{
		"address":  "http://localhost",
		"database": "",
		"username": "",
		"password": "",
	})
	assert.Nil(conn)
	assert.Error(err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	conn, err = Connect(map[string]interface{}{
		"address":  srv.URL,
		"database": "",
		"username": "",
		"password": "",
	})

	assert.NotNil(conn)
	assert.NoError(err)
}

func TestAddPoint(t *testing.T) {
	assert := assert.New(t)

	// Test add Point without tags
	connection := &Connection{
		config: map[string]interface{}{},
		points: make(chan *client.Point, 1),
	}

	connection.addPoint("name", models.Tags{}, models.Fields{"clients.total": 10}, time.Now())
	point := <-connection.points
	assert.NotNil(point)
	tags := point.Tags()
	assert.NotNil(tags)
	assert.NotEqual(tags["testtag2"], "value")

	// Test add Point with tags
	connection.config["tags"] = map[string]interface{}{
		"testtag": "value",
	}

	connection.addPoint("name", models.Tags{}, models.Fields{"clients.total": 10}, time.Now())
	point = <-connection.points
	assert.NotNil(point)
	tags = point.Tags()
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
