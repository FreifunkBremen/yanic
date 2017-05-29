package socket

import (
	"encoding/json"
	"log"
	"net"
)

func (conn *Connection) handleSocketConnection(ln net.Listener) {
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Println("[socket-database] error during connection of a client", err)
			continue
		}
		conn.clientMux.Lock()
		conn.clients[c.RemoteAddr()] = c
		conn.clientMux.Unlock()
	}
}

func (conn *Connection) sendJSON(msg Message) {
	conn.clientMux.Lock()
	for addr, c := range conn.clients {
		d := json.NewEncoder(c)

		err := d.Encode(&msg)
		if err != nil {
			log.Println("[socket-database] client has not recieve event:", err)
			c.Close()
			delete(conn.clients, addr)
		}
	}
	conn.clientMux.Unlock()
}
