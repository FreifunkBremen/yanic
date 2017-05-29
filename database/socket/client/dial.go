package yanic

import (
	"encoding/json"
	"log"
	"net"

	"github.com/FreifunkBremen/yanic/database/socket"
	"github.com/FreifunkBremen/yanic/runtime"
)

type Dialer struct {
	conn              net.Conn
	queue             chan socket.Message
	quit              chan struct{}
	NodeHandler       func(*runtime.Node)
	GlobalsHandler    func(*runtime.GlobalStats)
	PruneNodesHandler func()
}

func Dial(ctype, addr string) *Dialer {
	conn, err := net.Dial(ctype, addr)
	if err != nil {
		log.Panicf("yanic dial to %s:%s failed", ctype, addr)
	}
	dialer := &Dialer{
		conn:  conn,
		queue: make(chan socket.Message),
		quit:  make(chan struct{}),
	}

	return dialer
}

func (d *Dialer) Start() {
	go d.reciever()
	d.parser()
}
func (d *Dialer) Close() {
	d.conn.Close()
	close(d.quit)
}

func (d *Dialer) reciever() {
	decoder := json.NewDecoder(d.conn)
	var msg socket.Message

	for {
		select {
		case <-d.quit:
			close(d.queue)
			return
		default:
			decoder.Decode(&msg)
			d.queue <- msg
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
		case socket.MessageEventInsertGlobals:
			if d.GlobalsHandler != nil {
				var globals runtime.GlobalStats

				obj, _ := json.Marshal(msg.Body)
				json.Unmarshal(obj, &globals)

				d.GlobalsHandler(&globals)
			}
		case socket.MessageEventPruneNodes:
			if d.PruneNodesHandler != nil {
				d.PruneNodesHandler()
			}
		}
	}
}
