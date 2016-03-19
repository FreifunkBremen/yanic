package respond

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"log"
	"net"
	"reflect"
	"time"
)

//Collector for a specificle respond messages
type Collector struct {
	CollectType string
	connection  *net.UDPConn   // UDP socket
	queue       chan *Response // received responses
	onReceive   OnReceive
	msgType     reflect.Type

	// Ticker and stopper
	ticker *time.Ticker
	stop   chan interface{}
}

type OnReceive func(net.UDPAddr, interface{})

//NewCollector creates a Collector struct
func NewCollector(CollectType string, initialDelay time.Duration, interval time.Duration, msgStruct interface{}, onReceive OnReceive) *Collector {
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
		ticker:      time.NewTicker(interval),
		stop:        make(chan interface{}, 1),
		msgType:     reflect.TypeOf(msgStruct),
		onReceive:   onReceive,
	}

	go collector.receiver()
	go collector.parser()

	// Run senders
	go func() {
		time.Sleep(initialDelay)
		collector.sendOnce() // immediately
		collector.sender()   // periodically
	}()

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
	log.Println("request", coll.CollectType)
}

func (coll *Collector) sendPacket(address string) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Panic(err)
	}

	if _, err := coll.connection.WriteToUDP([]byte("GET "+coll.CollectType), addr); err != nil {
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
		if err := coll.parse(obj); err != nil {
			log.Println("unable to decode response from", obj.Address.String(), err, "\n", string(obj.Raw))
		}
	}
}

func (coll *Collector) parse(response *Response) (err error) {
	// create new struct instance
	data := reflect.New(coll.msgType).Interface()

	// deflater
	reader := flate.NewReader(bytes.NewReader(response.Raw))
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(data); err == nil {
		coll.onReceive(response.Address, data)
	}

	return
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
