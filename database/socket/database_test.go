package socket

import (
	"encoding/json"
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
	config["address"] = "/tmp/yanic-database2.socket"

	conn, err := Connect(config)
	assert.NoError(err, "connection should work")
	assert.NotNil(conn)

	client, err := net.Dial("unix", "/tmp/yanic-database2.socket")
	assert.NoError(err, "connection should work")
	assert.NotNil(client)
	time.Sleep(time.Duration(3) * time.Microsecond)

	decoder := json.NewDecoder(client)
	var msg Message

	conn.InsertNode(&runtime.Node{})
	decoder.Decode(&msg)
	assert.Equal("insert_node", msg.Event)

	conn.InsertGlobals(&runtime.GlobalStats{}, time.Now())
	decoder.Decode(&msg)
	assert.Equal("insert_globals", msg.Event)

	conn.PruneNodes(time.Hour * 24 * 7)
	decoder.Decode(&msg)
	assert.Equal("prune_nodes", msg.Event)
	time.Sleep(time.Duration(3) * time.Microsecond)

	// to reach in sendJSON removing of disconnection
	conn.Close()

	conn.InsertNode(&runtime.Node{})
	err = decoder.Decode(&msg)
	assert.Error(err)

}
