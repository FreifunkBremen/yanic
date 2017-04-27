package socket

import (
	"encoding/json"
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
			continue
		}
		config.clients[c.RemoteAddr()] = c
	}
}

func (conn *Connection) sendJSON(msg EventMessage) {
	for i, c := range conn.clients {
		d := json.NewEncoder(c)

		err := d.Encode(&msg)
		if err != nil {
			c.Close()
			delete(conn.clients, i)
		}
	}
}
