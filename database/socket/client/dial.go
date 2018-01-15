package yanic

import (
	"encoding/json"
	"io"
	"log"
	"net"

	"github.com/FreifunkBremen/yanic/database/socket"
	"github.com/FreifunkBremen/yanic/runtime"
)

var queueMaxSize = 1024

type Dialer struct {
	conn              net.Conn
	queue             chan socket.Message
	quit              chan struct{}
	NodeHandler       func(*runtime.Node)
	LinkHandler       func(*runtime.Link)
	GlobalsHandler    func(*runtime.GlobalStats, string)
	PruneNodesHandler func()
}

func Dial(ctype, addr string) *Dialer {
	conn, err := net.Dial(ctype, addr)
	if err != nil {
		log.Panicf("[yanic-client] dial to %s:%s failed", ctype, addr)
	}
	dialer := &Dialer{
		conn:  conn,
		queue: make(chan socket.Message, queueMaxSize),
		quit:  make(chan struct{}),
	}

	return dialer
}

func (d *Dialer) Start() {
	go d.receiver()
	d.parser()
}
func (d *Dialer) Close() {
	d.conn.Close()
	close(d.quit)
}

func (d *Dialer) receiver() {
	decoder := json.NewDecoder(d.conn)
	var msg socket.Message

	for {
		select {
		case <-d.quit:
			close(d.queue)
			return
		default:
			err := decoder.Decode(&msg)
			if err != nil {
				if err == io.EOF {
					log.Printf("[yanic-client] connection closed: %s", err)
					d.conn.Close()
					close(d.quit)
					close(d.queue)
					return
				}
				log.Printf("[yanic-client] could not decode message: %s", err)
				continue
			}
			select {
			case d.queue <- msg:
			default:
				log.Println("[yanic-client] full queue, drop latest entry")
				<-d.queue
			}
		}
	}
}

func (d *Dialer) parser() {
	for msg := range d.queue {
		switch msg.Event {
		case socket.MessageEventInsertNode:
			if d.NodeHandler != nil {
				var node runtime.Node

				obj, _ := json.Marshal(msg.Body)
				json.Unmarshal(obj, &node)
				d.NodeHandler(&node)
			}
		case socket.MessageEventInsertLink:
			if d.GlobalsHandler != nil {
				var link runtime.Link

				obj, _ := json.Marshal(msg.Body)
				json.Unmarshal(obj, &link)

				d.LinkHandler(&link)
			}
		case socket.MessageEventInsertGlobals:
			if d.GlobalsHandler != nil {
				var globals runtime.GlobalStats

				obj, _ := json.Marshal(msg.Body)
				json.Unmarshal(obj, &globals)

				d.GlobalsHandler(&globals, msg.Site)
			}
		case socket.MessageEventPruneNodes:
			if d.PruneNodesHandler != nil {
				d.PruneNodesHandler()
			}
		}
	}
}
