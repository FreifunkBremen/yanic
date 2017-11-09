package yanic

import (
	"testing"
	"time"

	"github.com/FreifunkBremen/yanic/database/socket"
	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/stretchr/testify/assert"
)

func TestConnectError(t *testing.T) {
	assert := assert.New(t)
	assert.Panics(func() {
		Dial("tcp4", "127.0.0.1:30303")
	}, "could connect")
}

func TestReceiveMessages(t *testing.T) {
	assert := assert.New(t)
	server, err := socket.Connect(map[string]interface{}{
		"enable":  true,
		"type":    "tcp4",
		"address": "127.0.0.1:10339",
	})
	assert.NoError(err)

	// test for drop queue
	queueMaxSize = 1

	d := Dial("tcp4", "127.0.0.1:10339")
	assert.NotNil(d)
	go d.Start()
	time.Sleep(5 * time.Millisecond)

	executed := false
	d.NodeHandler = func(node *runtime.Node) {
		executed = true
	}
	server.InsertNode(nil)
	time.Sleep(5 * time.Millisecond)
	assert.True(executed, "node not inserted")

	executed = false
	d.GlobalsHandler = func(stats *runtime.GlobalStats) {
		executed = true
	}
	server.InsertGlobals(nil, time.Now())
	time.Sleep(5 * time.Millisecond)
	assert.True(executed, "global stats not inserted")

	executed = false
	d.PruneNodesHandler = func() {
		executed = true
	}
	server.PruneNodes(time.Hour * 24 * 7)
	time.Sleep(5 * time.Millisecond)
	assert.True(executed, "node not pruned")

	// test for drop queue (only visible at test coverage)
	server.InsertNode(&runtime.Node{})
	server.InsertNode(&runtime.Node{})
	server.InsertNode(&runtime.Node{})
	time.Sleep(5 * time.Millisecond)

	d.Close()

	time.Sleep(5 * time.Millisecond)
	executed = false
	server.InsertNode(&runtime.Node{})
	time.Sleep(5 * time.Millisecond)
	assert.False(executed, "message re")

	server.Close()
}
