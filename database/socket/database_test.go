package socket

import (
	"net"
	"testing"
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestStartup(t *testing.T) {
	assert := assert.New(t)

	config := make(map[string]interface{})

	config["enable"] = false
	conn, err := Connect(config)
	assert.Nil(conn)

	config["enable"] = true
	config["type"] = ""
	config["address"] = ""
	conn, err = Connect(config)
	assert.Error(err, "connection should not work")
	assert.Nil(conn)

	config["type"] = "unix"
	config["address"] = "/tmp/yanic-database.socket"

	conn, err = Connect(config)
	assert.NoError(err, "connection should work")
	assert.NotNil(conn)

	conn.Close()
}

func TestClient(t *testing.T) {
	assert := assert.New(t)

	config := make(map[string]interface{})

	config["enable"] = true
	config["type"] = "unix"
	config["address"] = "/tmp/yanic-database.socket"

	conn, err := Connect(config)
	assert.NoError(err, "connection should work")
	assert.NotNil(conn)

	client, err := net.Dial("unix", "/tmp/yanic-database.socket")
	assert.NoError(err, "connection should work")
	assert.NotNil(client)
	time.Sleep(1)

	conn.InsertNode(&runtime.Node{})
	conn.InsertGlobals(&runtime.GlobalStats{}, time.Now())
	conn.PruneNodes(time.Duration(3))

	err = client.Close()
	assert.NoError(err, "disconnect should work")
	time.Sleep(1)
	conn.InsertNode(&runtime.Node{})

	conn.Close()
}
