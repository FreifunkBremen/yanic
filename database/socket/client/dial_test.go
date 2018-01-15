package yanic

import (
	"log"
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

func TestReceiveMessages(t *testing.T) {
	assert := assert.New(t)
	server, err := socket.Connect(map[string]interface{}{
		"type":    "tcp4",
		"address": "127.0.0.1:10339",
	})
	assert.NoError(err)
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}

	// test for drop queue
	queueMaxSize = 1

	d := Dial("tcp4", "127.0.0.1:10339")
	d.NodeHandler = func(node *runtime.Node) {
		wg.Done()
	}
	d.LinkHandler = func(link *runtime.Link) {
		wg.Done()
	}
	d.GlobalsHandler = func(stats *runtime.GlobalStats, site string) {
		wg.Done()
	}
	d.PruneNodesHandler = func() {
		wg.Done()
	}
	assert.NotNil(d)
	go d.Start()
	time.Sleep(5 * time.Millisecond)

	wg.Add(1)
	server.InsertNode(&runtime.Node{})
	log.Print("[run] wait for insert node")
	wg.Wait()
	log.Print("[result] node inserted")

	wg.Add(1)
	server.InsertLink(nil, time.Now())
	log.Print("[run] wait for insert link")
	wg.Wait()
	log.Print("[result] link inserted")

	wg.Add(1)
	server.InsertGlobals(nil, time.Now(), "global")
	log.Print("[run] wait for insert globals")
	wg.Wait()
	log.Print("[result] global stats inserted")

	wg.Add(1)
	server.PruneNodes(time.Hour * 24 * 7)
	log.Print("[run] wait for prune node")
	wg.Wait()
	log.Print("[result] node pruned")

	//TODO test query overload

	d.Close()
	log.Print("closed connection")

	server.InsertNode(&runtime.Node{})
	log.Print("handle closed client")

	server.Close()
}
