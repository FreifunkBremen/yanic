package main

import (
	"encoding/json"
	"log"
	"net"
	"reflect"
	"strings"
	"time"
)

const (
	maxDatagramSize = 8192
)

type Collector struct {
	collectType string
	connection  *net.UDPConn // UDP socket
	queue       chan string  // received responses
}

func NewCollector(collectType string) *Collector {
	// Parse address
	addr, err := net.ResolveUDPAddr("udp", "[::]:0")
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
		collectType: collectType,
		connection:  conn,
		queue:       make(chan string, 100),
	}

	go collector.sendOnce()
	go collector.sender()

	go collector.receiver()
	go collector.parser()

	return collector
}

func (coll *Collector) Close() {
	coll.connection.Close()
	close(coll.queue)
}

func (coll *Collector) sendOnce() {
	coll.sendPacket("[2a06:8782:ffbb:1337:c24a:ff:fe2c:c7ac]:1001")
	coll.sendPacket("[2001:bf7:540:0:32b5:c2ff:fe6e:99d5]:1001")
}

func (coll *Collector) sendPacket(address string) {
	addr, err := net.ResolveUDPAddr("udp", address)
	check(err)

	coll.connection.WriteToUDP([]byte(coll.collectType), addr)
}

func (coll *Collector) sender() {
	c := time.Tick(collectInterval)

	for range c {
		coll.sendOnce()
	}
}

func (coll *Collector) parser() {
	for str := range coll.queue {
		coll.parseSingle(str)
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

	node := nodes.get(nodeId)

	// Set result
	elem := reflect.ValueOf(node).Elem()
	field := elem.FieldByName(strings.Title(coll.collectType))
	field.Set(reflect.ValueOf(result))
}

func (coll *Collector) receiver() {
	b := make([]byte, maxDatagramSize)
	for {
		n, src, err := coll.connection.ReadFromUDP(b)

		if err != nil {
			log.Println("ReadFromUDP failed:", err)
			return
		}
		coll.queue <- string(b[:n])
		log.Println("received", coll.collectType, "from", src)
	}
}
