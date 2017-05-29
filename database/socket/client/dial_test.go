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
		Dial("tcp6", "[::]:30303")
	}, "could connect")
}

func TestRecieveMessages(t *testing.T) {
	assert := assert.New(t)
	server, err := socket.Connect(map[string]interface{}{
		"enable":  true,
		"type":    "tcp6",
		"address": "[::]:1337",
	})
	assert.NoError(err)

	d := Dial("tcp6", "[::]:1337")
	assert.NotNil(d)
	go d.Start()
	time.Sleep(5 * time.Millisecond)

	runned := false
	d.NodeHandler = func(node *runtime.Node) {
		runned = true
	}
	server.InsertNode(nil)
	time.Sleep(5 * time.Millisecond)
	assert.True(runned, "node not inserted")

	runned = false
	d.GlobalsHandler = func(stats *runtime.GlobalStats) {
		runned = true
	}
	server.InsertGlobals(nil, time.Now())
	time.Sleep(5 * time.Millisecond)
	assert.True(runned, "global stats not inserted")

	runned = false
	d.PruneNodesHandler = func() {
		runned = true
	}
	server.PruneNodes(time.Hour * 24 * 7)
	time.Sleep(5 * time.Millisecond)
	assert.True(runned, "node not pruned")

	d.Close()

	time.Sleep(5 * time.Millisecond)
	runned = false
	d.PruneNodesHandler = func() {
		runned = true
	}
	server.PruneNodes(time.Hour * 24 * 7)
	time.Sleep(5 * time.Millisecond)
	assert.False(runned, "message recieve")

	server.Close()
}
