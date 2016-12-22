package respond

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"reflect"
	"time"

	"github.com/FreifunkBremen/respond-collector/data"
	"github.com/FreifunkBremen/respond-collector/database"
	"github.com/FreifunkBremen/respond-collector/models"
)

//Collector for a specificle respond messages
type Collector struct {
	CollectType   string
	connection    *net.UDPConn   // UDP socket
	queue         chan *Response // received responses
	msgType       reflect.Type
	multicastAddr string
	db            *database.DB
	nodes         *models.Nodes
	// Ticker and stopper
	ticker *time.Ticker
	stop   chan interface{}
}

// Creates a Collector struct
func NewCollector(db *database.DB, nodes *models.Nodes, iface string) *Collector {
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
		connection:    conn,
		db:            db,
		nodes:         nodes,
		multicastAddr: net.JoinHostPort(multiCastGroup+"%"+iface, port),
		queue:         make(chan *Response, 400),
		stop:          make(chan interface{}, 1),
	}

	go collector.receiver()
	go collector.parser()

	if collector.db != nil {
		go collector.globalStatsWorker()
	}

	return collector
}

// Start Collector
func (coll *Collector) Start(interval time.Duration) {
	if coll.ticker != nil {
		panic("already started")
	}

	coll.ticker = time.NewTicker(interval)
	go func() {
		coll.sendOnce() // immediately
		coll.sender()   // periodically
	}()
}

// Close Collector
func (coll *Collector) Close() {
	// stop ticker
	if coll.ticker != nil {
		coll.ticker.Stop()
		close(coll.stop)
	}

	coll.connection.Close()
	close(coll.queue)
}

func (coll *Collector) sendOnce() {
	coll.SendPacket(coll.multicastAddr)
}

// Sends a UDP request to the given unicast or multicast address
func (coll *Collector) SendPacket(address string) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Panic(err)
	}

	if _, err := coll.connection.WriteToUDP([]byte("GET nodeinfo statistics neighbours"), addr); err != nil {
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
			// send the multicast packet to request per-node statistics
			coll.sendOnce()
		}
	}
}

func (coll *Collector) parser() {
	for obj := range coll.queue {
		if data, err := obj.parse(); err != nil {
			log.Println("unable to decode response from", obj.Address.String(), err, "\n", string(obj.Raw))
		} else {
			coll.saveResponse(obj.Address, data)
		}
	}
}

func (res *Response) parse() (*data.ResponseData, error) {
	// Deflate
	deflater := flate.NewReader(bytes.NewReader(res.Raw))
	defer deflater.Close()

	// Unmarshal
	rdata := &data.ResponseData{}
	err := json.NewDecoder(deflater).Decode(rdata)

	return rdata, err
}

func (coll *Collector) saveResponse(addr net.UDPAddr, res *data.ResponseData) {
	// Search for NodeID
	var nodeId string
	if val := res.NodeInfo; val != nil {
		nodeId = val.NodeId
	} else if val := res.Neighbours; val != nil {
		nodeId = val.NodeId
	} else if val := res.Statistics; val != nil {
		nodeId = val.NodeId
	}

	// Updates nodes if NodeID found
	if len(nodeId) != 12 {
		log.Printf("invalid NodeID '%s' from %s", nodeId, addr.String())
		return
	}
	node := coll.nodes.Update(nodeId, res)

	if coll.db != nil && node.Statistics != nil {
		coll.db.Add(nodeId, node)
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

func (coll *Collector) globalStatsWorker() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-coll.stop:
			return
		case <-ticker.C:
			coll.saveGlobalStats()
		}
	}
}

// saves global statistics
func (coll *Collector) saveGlobalStats() {
	stats := models.NewGlobalStats(coll.nodes)

	coll.db.AddPoint(database.MeasurementGlobal, nil, stats.Fields(), time.Now())
	coll.db.AddCounterMap(database.MeasurementFirmware, stats.Firmwares)
	coll.db.AddCounterMap(database.MeasurementModel, stats.Models)
}
