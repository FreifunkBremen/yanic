package socket

import (
	"encoding/json"
	"log"
	"net"
	"strings"
)

func (conn *Connection) handleSocketConnection(ln net.Listener) {
	for {
		c, err := ln.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Println("[socket-database] connection already closed, no new client")
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
	for msg := range conn.buffer {
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
	conn.buffer <- msg
}
