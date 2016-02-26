package websocketserver

import (
	"fmt"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

var maxID = 0

//Client struct
type Client struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan interface{}
	doneCh chan bool
}

//NewClient creates a new Client
func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxID++

	return &Client{
		id:     maxID,
		ws:     ws,
		server: server,
		ch:     make(chan interface{}, channelBufSize),
		doneCh: make(chan bool),
	}
}

//GetConnection the websocket connection of a listen client
func (c *Client) GetConnection() *websocket.Conn {
	return c.ws
}

//Write send the msg informations to the clients
func (c *Client) Write(msg interface{}) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("Client %d is disconnected.", c.id)
		c.server.Err(err)
	}
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	c.listen()
}

// listen for new msg informations
func (c *Client) listen() {
	for {
		select {
		case msg := <-c.ch:
			err := websocket.JSON.Send(c.ws, msg)
			if err != nil {
				c.doneCh <- true
			}
		case gone := <-c.doneCh:
			if gone {
				c.server.Del(c)
				err := fmt.Errorf("Client %d is disconnected.", c.id)
				c.server.Err(err)
			}
		}
	}
}
