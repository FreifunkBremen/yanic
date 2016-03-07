package respond

import (
	"log"
	"net"
	"time"
)

const (
	// default multicast group used by announced
	multiCastGroup string = "ff02:0:0:0:0:0:2:1001"

	// default udp port used by announced
	port string = "1001"

	// maximum receivable size
	maxDataGramSize int = 8192
)

//Response of the respond request
type Response struct {
	Address net.UDPAddr
	Raw     []byte
}

//Collector for a specificle respond messages
type Collector struct {
	CollectType string
	connection  *net.UDPConn   // UDP socket
	queue       chan *Response // received responses
	parse       func(coll *Collector, res *Response)
}

//NewCollector creates a Collector struct
func NewCollector(CollectType string, parseFunc func(coll *Collector, res *Response)) *Collector {
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
	}

	go collector.receiver()
	go collector.parser()

	collector.sendOnce()

	return collector
}

//Close Collector
func (coll *Collector) Close() {
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

func (coll *Collector) sender(collectInterval time.Duration) {
	c := time.Tick(collectInterval)

	for range c {
		// TODO break condition
		coll.sendOnce()
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
		log.Println("received", coll.CollectType, "from", src)
	}
}
