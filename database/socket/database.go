package socket

/*
 * This socket database is to run another service
 * (without flooding the network with respondd packages)
 * e.g. https://github.com/FreifunkBremen/freifunkmanager
 */

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/runtime"
)

var queueMaxSize = 1024

type Connection struct {
	database.Connection
	listener  net.Listener
	clients   map[net.Addr]net.Conn
	clientMux sync.Mutex
	queue     chan *Message
	running   sync.WaitGroup
}

func init() {
	database.RegisterAdapter("socket", Connect)
}

func Connect(config map[string]interface{}) (database.Connection, error) {
	ln, err := net.Listen(config["type"].(string), config["address"].(string))
	if err != nil {
		return nil, err
	}
	conn := &Connection{
		listener: ln,
		clients:  make(map[net.Addr]net.Conn),
		queue:    make(chan *Message, queueMaxSize),
	}
	go conn.handleSocketConnection(ln)
	go conn.writer()

	log.Println("[socket-database] listen on: ", ln.Addr())

	return conn, nil
}

func (conn *Connection) InsertNode(node *runtime.Node) {
	conn.sendJSON(&Message{Event: MessageEventInsertNode, Body: node})
}

func (conn *Connection) InsertLink(link *runtime.Link, time time.Time) {
	conn.sendJSON(&Message{Event: MessageEventInsertLink, Body: link})
}

func (conn *Connection) InsertGlobals(stats *runtime.GlobalStats, time time.Time, site string) {
	conn.sendJSON(&Message{Event: MessageEventInsertGlobals, Body: stats, Site: site})
}

func (conn *Connection) PruneNodes(deleteAfter time.Duration) {
	conn.sendJSON(&Message{Event: MessageEventPruneNodes})
}

func (conn *Connection) Close() {
	conn.clientMux.Lock()
	for _, c := range conn.clients {
		c.Close()
	}
	conn.clientMux.Unlock()
	conn.listener.Close()
	conn.running.Wait()
}
