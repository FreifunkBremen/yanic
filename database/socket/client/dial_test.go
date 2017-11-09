package yanic

import (
	"sync"
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

type SafeBoolean struct {
	value bool
	sync.Mutex
}

func (s *SafeBoolean) Set(value bool) {
	s.Lock()
	s.value = value
	s.Unlock()
}

func (s *SafeBoolean) Get() bool {
	s.Lock()
	defer s.Unlock()
	return s.value
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
	executed := SafeBoolean{value: false}
	d.NodeHandler = func(node *runtime.Node) {
		executed.Set(true)
	}
	d.LinkHandler = func(link *runtime.Link) {
		executed.Set(true)
	}
	d.GlobalsHandler = func(stats *runtime.GlobalStats) {
		executed.Set(true)
	}
	d.PruneNodesHandler = func() {
		executed.Set(true)
	}
	assert.NotNil(d)
	go d.Start()
	time.Sleep(5 * time.Millisecond)

	server.InsertNode(nil)
	time.Sleep(5 * time.Millisecond)
	assert.True(executed.Get(), "node not inserted")

	executed.Set(false)
	server.InsertLink(nil, time.Now())
	time.Sleep(5 * time.Millisecond)
	assert.True(executed.Get(), "link not inserted")

	executed.Set(false)
	server.InsertGlobals(nil, time.Now())
	time.Sleep(5 * time.Millisecond)
	assert.True(executed.Get(), "global stats not inserted")

	executed.Set(false)
	server.PruneNodes(time.Hour * 24 * 7)
	time.Sleep(5 * time.Millisecond)
	assert.True(executed.Get(), "node not pruned")

	// test for drop queue (only visible at test coverage)
	server.InsertNode(&runtime.Node{})
	server.InsertNode(&runtime.Node{})
	server.InsertNode(&runtime.Node{})
	time.Sleep(5 * time.Millisecond)

	d.Close()

	time.Sleep(5 * time.Millisecond)
	executed.Set(false)
	server.InsertNode(&runtime.Node{})
	time.Sleep(5 * time.Millisecond)
	assert.False(executed.Get(), "message re")

	server.Close()
}
