package websocketserver

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

//Server struct
type Server struct {
	pattern   string
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *struct{}
	closeCh   chan bool
	errCh     chan error
}

//NewServer creates a new node server
func NewServer(pattern string) *Server {
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *struct{})
	closeCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		clients,
		addCh,
		delCh,
		sendAllCh,
		closeCh,
		errCh,
	}
}

//Add a node listen client
func (s *Server) Add(c *Client) {
	s.addCh <- c
}

//Del a node listen client
func (s *Server) Del(c *Client) {
	s.delCh <- c
}

//SendAll to all listen clients refreshed information of a node
func (s *Server) SendAll(node *struct{}) {
	s.sendAllCh <- node
}

//Close stops node server
func (s *Server) Close() {
	s.closeCh <- true
}

//Err send to server
func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendAll(node *struct{}) {
	for _, c := range s.clients {
		c.Write(node)
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {

	log.Println("Listening Server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
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
			s.sendAll(node)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.closeCh:
			return
		}
	}
}
