package main

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

const (
	maxDatagramSize = 8192
)

type Node struct {
	Firstseen  time.Time   `json:"firstseen"`
	Lastseen   time.Time   `json:"lastseen"`
	Statistics interface{} `json:"statistics"`
	Nodeinfo   interface{} `json:"nodeinfo"`
}

type Collector struct {
	connection *net.UDPConn     // UDP socket
	queue      chan string      // received responses
	nodes      map[string]*Node // the current nodemap
}

func NewCollector() *Collector {
	// Parse address
	addr, err := net.ResolveUDPAddr("udp", "[::]:1001")
	if err != nil {
		log.Panic(err)
	}

	// Open socket
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Panic(err)
	}
	conn.SetReadBuffer(maxDatagramSize)

	collector := &Collector{
		connection: conn,
		queue:      make(chan string, 100),
		nodes:      make(map[string]*Node),
	}

	go collector.receiver()
	go collector.parser()

	return collector
}

func (coll *Collector) Close() {
	coll.connection.Close()
	close(coll.queue)
}

func (coll *Collector) send(address string) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Panic(err)
	}
	coll.connection.WriteToUDP([]byte("nodeinfo"), addr)
}

func (coll *Collector) print() {
	b, err := json.Marshal(coll.nodes)
	if err != nil {
		log.Panic(err)
	}
	log.Println(string(b))
}

func (coll *Collector) parser() {
	for str := range coll.queue {
		coll.parseSingle(str)
		coll.print()
	}
}

// Parst die Rückgabe
func (coll *Collector) parseSingle(str string) {
	var result map[string]interface{}
	json.Unmarshal([]byte(str), &result)

	nodeId, _ := result["node_id"].(string)

	if nodeId == "" {
		log.Println("unable to parse node_id")
		return
	}

	now := time.Now()
	node, _ := coll.nodes[nodeId]

	if node == nil {
		node = &Node{
			Firstseen: now,
		}
		coll.nodes[nodeId] = node
	}

	node.Lastseen = now
	node.Nodeinfo = result
}

func (coll *Collector) receiver() {
	b := make([]byte, maxDatagramSize)
	for {
		n, _, err := coll.connection.ReadFromUDP(b)

		if err != nil {
			log.Println("ReadFromUDP failed:", err)
		} else {
			coll.queue <- string(b[:n])
		}
	}
}
