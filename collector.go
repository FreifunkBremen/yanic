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
	// default multicast group used by announced
	MultiCastGroup string = "ff02:0:0:0:0:0:2:1001"

	// default udp port used by announced
	Port string = "1001"

	// maximum receivable size
	MaxDataGramSize int = 8192
)

type Response struct {
	Address net.UDPAddr
	Raw     []byte
}

type Collector struct {
	collectType string
	connection  *net.UDPConn   // UDP socket
	queue       chan *Response // received responses
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
	conn.SetReadBuffer(MaxDataGramSize)

	collector := &Collector{
		collectType: collectType,
		connection:  conn,
		queue:       make(chan *Response, 400),
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
	coll.sendPacket(net.JoinHostPort(MultiCastGroup,Port))
}

func (coll *Collector) sendPacket(address string) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Panic(err)
	}

	coll.connection.WriteToUDP([]byte(coll.collectType), addr)
}

func (coll *Collector) sender() {
	c := time.Tick(collectInterval)

	for range c {
		// TODO break condition
		coll.sendOnce()
	}
}

func (coll *Collector) parser() {
	for obj := range coll.queue {
		coll.parse(obj)
	}
}

// Parses a response
func (coll *Collector) parse(res *Response) {
	var result map[string]interface{}
	json.Unmarshal(res.Raw, &result)

	nodeId, _ := result["node_id"].(string)

	if nodeId == "" {
		log.Println("unable to parse node_id")
		return
	}

	node := nodes.get(nodeId)

	// Set result
	elem := reflect.ValueOf(node).Elem()
	field := elem.FieldByName(strings.Title(coll.collectType))

	log.Println(field)
	log.Println(result)

	if !reflect.DeepEqual(field,result){
		nodeserver.SendAll(node)
	}

	field.Set(reflect.ValueOf(result))
}

func (coll *Collector) receiver() {
	buf := make([]byte, MaxDataGramSize)
	for {
		n, src, err := coll.connection.ReadFromUDP(buf)

		if err != nil {
			log.Println("ReadFromUDP failed:", err)
			return
		}

		raw := make([]byte, n)
		copy(raw, buf)

		coll.queue <- &Response{
			Address: *src,
			Raw:     raw,
		}
		log.Println("received", coll.collectType, "from", src)
	}
}
