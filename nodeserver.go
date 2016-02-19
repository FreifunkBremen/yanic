package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// Node server.
type NodeServer struct {
	pattern   string
	clients   map[int]*NodeClient
	addCh     chan *NodeClient
	delCh     chan *NodeClient
	sendAllCh chan *Node
	doneCh    chan bool
	errCh     chan error
}

// Create new node server.
func NewNodeServer(pattern string) *NodeServer {
	clients := make(map[int]*NodeClient)
	addCh := make(chan *NodeClient)
	delCh := make(chan *NodeClient)
	sendAllCh := make(chan *Node)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &NodeServer{
		pattern,
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}

func (s *NodeServer) Add(c *NodeClient) {
	s.addCh <- c
}

func (s *NodeServer) Del(c *NodeClient) {
	s.delCh <- c
}

func (s *NodeServer) SendAll(node *Node) {
	s.sendAllCh <- node
}

func (s *NodeServer) Done() {
	s.doneCh <- true
}

func (s *NodeServer) Err(err error) {
	s.errCh <- err
}

func (s *NodeServer) sendAll(node *Node) {
	for _, c := range s.clients {
		c.Write(node)
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *NodeServer) Listen() {

	log.Println("Listening NodeServer...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewNodeClient(ws, s)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)

		// broadcast message for all clients
		case node := <-s.sendAllCh:
			log.Println("Send all:", node)
			s.sendAll(node)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
