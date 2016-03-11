package respond

import (
	"log"
	"net"
	"time"
)

//Collector for a specificle respond messages
type Collector struct {
	CollectType string
	connection  *net.UDPConn   // UDP socket
	queue       chan *Response // received responses
	parse       func(coll *Collector, res *Response)

	// Ticker and stopper
	ticker *time.Ticker
	stop   chan interface{}
}

type ParseFunc func(coll *Collector, res *Response)

//NewCollector creates a Collector struct
func NewCollector(CollectType string, interval time.Duration, parseFunc ParseFunc) *Collector {
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
	conn.SetReadBuffer(maxDataGramSize)

	collector := &Collector{
		CollectType: CollectType,
		connection:  conn,
		queue:       make(chan *Response, 400),
		parse:       parseFunc,
		ticker:      time.NewTicker(interval),
		stop:        make(chan interface{}, 1),
	}

	go collector.receiver()
	go collector.parser()
	go collector.sender()

	return collector
}

// Close Collector
func (coll *Collector) Close() {
	// stop ticker
	coll.ticker.Stop()
	coll.stop <- nil

	coll.connection.Close()
	close(coll.queue)
}

func (coll *Collector) sendOnce() {
	coll.sendPacket(net.JoinHostPort(multiCastGroup, port))
	log.Println("send request")
}

func (coll *Collector) sendPacket(address string) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Panic(err)
	}

	if _, err := coll.connection.WriteToUDP([]byte(coll.CollectType), addr); err != nil {
		log.Println("WriteToUDP failed:", err)
	}
}

// send packets continously
func (coll *Collector) sender() {
	for {
		select {
		case <-coll.stop:
			return
		case <-coll.ticker.C:
			coll.sendOnce()
		}
	}
}

func (coll *Collector) parser() {
	for obj := range coll.queue {
		coll.parse(coll, obj)
	}
}

func (coll *Collector) receiver() {
	buf := make([]byte, maxDataGramSize)
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
	}
}
