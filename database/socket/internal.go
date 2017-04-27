package socket

import (
	"encoding/json"
	"log"
	"net"
)

type EventMessage struct {
	Event string      `json:"event"`
	Body  interface{} `json:"body,omitempty"`
}

func (config *Connection) handleSocketConnection(ln net.Listener) {
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Println("[socket-database] error during connection of a client", err)
			continue
		}
		config.clients[c.RemoteAddr()] = c
	}
}

func (conn *Connection) sendJSON(msg EventMessage) {
	for addr, c := range conn.clients {
		d := json.NewEncoder(c)

		err := d.Encode(&msg)
		if err != nil {
			log.Println("[socket-database] client has not recieve event:", err)
			c.Close()
			delete(conn.clients, addr)
		}
	}
}
