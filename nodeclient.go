package main

import (
	"fmt"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

var maxId int = 0

// Node client.
type NodeClient struct {
	id     int
	ws     *websocket.Conn
	server *NodeServer
	ch     chan *Node
	doneCh chan bool
}

// Create new node client.
func NewNodeClient(ws *websocket.Conn, server *NodeServer) *NodeClient {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxId++
	ch := make(chan *Node, channelBufSize)
	doneCh := make(chan bool)

	return &NodeClient{maxId, ws, server, ch, doneCh}
}

func (c *NodeClient) Conn() *websocket.Conn {
	return c.ws
}

func (c *NodeClient) Write(node *Node) {
	select {
	case c.ch <- node:
	default:
		c.server.Del(c)
		err := fmt.Errorf("NodeClient %d is disconnected.", c.id)
		c.server.Err(err)
	}
}

func (c *NodeClient) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *NodeClient) Listen() {
	c.listenWrite()
}

// Listen write request via chanel
func (c *NodeClient) listenWrite() {
	for {
		select {

		// send message to the client
		case node := <-c.ch:
			websocket.JSON.Send(c.ws, node)

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}
