package socket

import (
	"encoding/json"
	"fmt"
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

	var msg Message

	fmt.Println("[run] insert node")
	conn.InsertNode(&runtime.Node{})
	err = json.NewDecoder(client).Decode(&msg)
	fmt.Println("[result] insert node")
	assert.NoError(err)
	assert.Equal(MessageEventInsertNode, msg.Event)

	fmt.Println("[run] insert globals")
	conn.InsertGlobals(&runtime.GlobalStats{}, time.Now(), runtime.GLOBAL_SITE)
	err = json.NewDecoder(client).Decode(&msg)
	fmt.Println("[result] insert globals")
	assert.NoError(err)
	assert.Equal(MessageEventInsertGlobals, msg.Event)

	fmt.Println("[run] insert link")
	conn.InsertLink(&runtime.Link{}, time.Now())
	err = json.NewDecoder(client).Decode(&msg)
	fmt.Println("[result] insert link")
	assert.NoError(err)
	assert.Equal(MessageEventInsertLink, msg.Event)

	fmt.Println("[run] prune nodes")
	conn.PruneNodes(time.Hour * 24 * 7)
	err = json.NewDecoder(client).Decode(&msg)
	fmt.Println("[result] prune nodes")
	assert.NoError(err)
	assert.Equal(MessageEventPruneNodes, msg.Event)

	//TODO test for drop queue (only visible at test coverage)

	// to reach in sendJSON removing of disconnection
	conn.Close()

}
