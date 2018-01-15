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

	config["type"] = ""
	config["address"] = ""
	conn, err := Connect(config)
	assert.Error(err, "connection should not work")
	assert.Nil(conn)

	config["type"] = "tcp4"
	config["address"] = "127.0.0.1:10337"

	conn, err = Connect(config)
	assert.NoError(err, "connection should work")
	assert.NotNil(conn)

	conn.Close()
}

func TestClient(t *testing.T) {
	assert := assert.New(t)

	config := make(map[string]interface{})

	config["enable"] = true
	config["type"] = "tcp4"
	config["address"] = "127.0.0.1:10338"

	// test for drop queue
	queueMaxSize = 1

	conn, err := Connect(config)
	assert.NoError(err, "connection should work")
	assert.NotNil(conn)

	client, err := net.Dial("tcp4", "127.0.0.1:10338")
	assert.NoError(err, "connection should work")
	assert.NotNil(client)
	time.Sleep(time.Duration(3) * time.Microsecond)

	decoder := json.NewDecoder(client)
	var msg Message

	conn.InsertNode(&runtime.Node{})
	err = decoder.Decode(&msg)
	assert.NoError(err)
	assert.Equal("insert_node", msg.Event)

	conn.InsertGlobals(&runtime.GlobalStats{}, time.Now(), "global")
	err = decoder.Decode(&msg)
	assert.NoError(err)
	assert.Equal("insert_globals", msg.Event)

	conn.InsertLink(&runtime.Link{}, time.Now())
	err = decoder.Decode(&msg)
	assert.NoError(err)
	assert.Equal("insert_link", msg.Event)

	conn.PruneNodes(time.Hour * 24 * 7)
	err = decoder.Decode(&msg)
	assert.NoError(err)
	assert.Equal("prune_nodes", msg.Event)

	// test for drop queue (only visible at test coverage)
	conn.InsertNode(&runtime.Node{})
	conn.InsertNode(&runtime.Node{})
	conn.InsertNode(&runtime.Node{})

	// to reach in sendJSON removing of disconnection
	conn.Close()

	conn.InsertNode(&runtime.Node{})
	err = decoder.Decode(&msg)
	assert.Error(err)

}
