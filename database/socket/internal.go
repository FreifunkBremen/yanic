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
	for i, c := range conn.clients {
		d := json.NewEncoder(c)

		err := d.Encode(&msg)
		if err != nil {
			err = c.Close()
			if err != nil {
				log.Println("[socket-database] connection could not close after error on sending event:", err)
			}
			delete(conn.clients, i)
		}
	}
}
