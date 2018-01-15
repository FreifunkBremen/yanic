package socket

import (
	"encoding/json"
	"log"
	"net"
	"strings"
)

func (conn *Connection) handleSocketConnection(ln net.Listener) {
	conn.running.Add(1)
	defer conn.running.Done()
	for {
		c, err := ln.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Println("[socket-database] connection already closed, no new client")
				close(conn.queue)
				return
			}
			log.Println("[socket-database] error during connection of a client", err)
			continue
		}
		conn.clientMux.Lock()
		conn.clients[c.RemoteAddr()] = c
		conn.clientMux.Unlock()
	}
}

func (conn *Connection) writer() {
	conn.running.Add(1)
	defer conn.running.Done()
	for msg := range conn.queue {
		conn.clientMux.Lock()
		for addr, c := range conn.clients {
			err := json.NewEncoder(c).Encode(&msg)
			if err != nil {
				log.Println("[socket-database] client has not receive event:", err)
				c.Close()
				delete(conn.clients, addr)
			}
		}
		conn.clientMux.Unlock()
	}
}

func (conn *Connection) sendJSON(msg *Message) {
	select {
	case conn.queue <- msg:
	default:
		log.Println("[socket-database] full queue, drop lates entry")
		<-conn.queue
	}
}
