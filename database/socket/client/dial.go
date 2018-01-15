package yanic

import (
	"encoding/json"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/FreifunkBremen/yanic/database/socket"
	"github.com/FreifunkBremen/yanic/runtime"
)

var queueMaxSize = 1024

type Dialer struct {
	conn              net.Conn
	queue             chan socket.Message
	wg                sync.WaitGroup
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
	}

	return dialer
}

func (d *Dialer) Start() {
	d.wg.Add(1)
	go d.parser()
	d.receiver()
	d.wg.Done()
}
func (d *Dialer) Close() {
	d.conn.Close()
	d.wg.Wait()
}

func (d *Dialer) receiver() {
	var msg socket.Message
	for {
		err := json.NewDecoder(d.conn).Decode(&msg)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("[yanic-client] connection closed: %s", err)
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

				err := remashal(msg.Body, &node)
				if err != nil {
					break
				}

				d.NodeHandler(&node)
			}
		case socket.MessageEventInsertLink:
			if d.GlobalsHandler != nil {
				var link runtime.Link

				err := remashal(msg.Body, &link)
				if err != nil {
					break
				}

				d.LinkHandler(&link)
			}
		case socket.MessageEventInsertGlobals:
			if d.GlobalsHandler != nil {
				var globals runtime.GlobalStats

				err := remashal(msg.Body, &globals)
				if err != nil {
					break
				}

				d.GlobalsHandler(&globals, msg.Site)
			}
		case socket.MessageEventPruneNodes:
			if d.PruneNodesHandler != nil {
				d.PruneNodesHandler()
			}

		default:
			log.Printf("[yanic-client] unknown message: %s", msg.Event)
		}
	}
	log.Println("[yanic-client] close")
}

func remashal(in, out interface{}) (err error) {
	obj, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = json.Unmarshal(obj, out)
	return err
}
