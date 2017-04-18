package respond

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/database"
	"github.com/FreifunkBremen/yanic/jsontime"
	"github.com/FreifunkBremen/yanic/runtime"
)

// Collector for a specificle respond messages
type Collector struct {
	connection *net.UDPConn   // UDP socket
	queue      chan *Response // received responses
	iface      string
	db         database.Connection
	nodes      *runtime.Nodes
	interval   time.Duration // Interval for multicast packets
	stop       chan interface{}
}

// NewCollector creates a Collector struct
func NewCollector(db database.Connection, nodes *runtime.Nodes, iface string, port int) *Collector {
	linkLocalAddr, err := getLinkLocalAddr(iface)
	if err != nil {
		log.Panic(err)
	}

	// Open socket
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   linkLocalAddr,
		Port: port,
		Zone: iface,
	})
	if err != nil {
		log.Panic(err)
	}
	conn.SetReadBuffer(maxDataGramSize)

	collector := &Collector{
		connection: conn,
		db:         db,
		nodes:      nodes,
		iface:      iface,
		queue:      make(chan *Response, 400),
		stop:       make(chan interface{}),
	}

	go collector.receiver()
	go collector.parser()

	if collector.db != nil {
		go collector.globalStatsWorker()
	}

	return collector
}

// Returns the first link local unicast address for the given interface name
func getLinkLocalAddr(ifname string) (net.IP, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, err
	}

	addresses, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addresses {
		if ipnet := addr.(*net.IPNet); ipnet.IP.IsLinkLocalUnicast() {
			return ipnet.IP, nil
		}
	}
	return nil, fmt.Errorf("unable to find link local unicast address for %s", ifname)
}

// Start Collector
func (coll *Collector) Start(interval time.Duration) {
	if coll.interval != 0 {
		panic("already started")
	}
	if interval <= 0 {
		panic("invalid collector interval")
	}
	coll.interval = interval

	go func() {
		coll.sendOnce() // immediately
		coll.sender()   // periodically
	}()
}

// Close Collector
func (coll *Collector) Close() {
	close(coll.stop)
	coll.connection.Close()
	close(coll.queue)
}

func (coll *Collector) sendOnce() {
	now := jsontime.Now()
	coll.sendMulticast()

	// Wait for the multicast responses to be processed and send unicasts
	time.Sleep(coll.interval / 2)
	coll.sendUnicasts(now)
}

func (coll *Collector) sendMulticast() {
	log.Println("sending multicast")
	coll.SendPacket(net.ParseIP(multiCastGroup))
}

// Send unicast packets to nodes that did not answer the multicast
func (coll *Collector) sendUnicasts(seenBefore jsontime.Time) {
	seenAfter := seenBefore.Add(-time.Minute * 10)

	// Select online nodes that has not been seen recently
	nodes := coll.nodes.Select(func(n *runtime.Node) bool {
		return n.Lastseen.After(seenAfter) && n.Lastseen.Before(seenBefore) && n.Address != nil
	})

	// Send unicast packets
	log.Printf("sending unicast to %d nodes", len(nodes))
	for _, node := range nodes {
		coll.SendPacket(node.Address)
		time.Sleep(10 * time.Millisecond)
	}
}

// SendPacket sends a UDP request to the given unicast or multicast address
func (coll *Collector) SendPacket(address net.IP) {
	addr := net.UDPAddr{
		IP:   address,
		Port: port,
		Zone: coll.iface,
	}

	if _, err := coll.connection.WriteToUDP([]byte("GET nodeinfo statistics neighbours"), &addr); err != nil {
		log.Println("WriteToUDP failed:", err)
	}
}

// send packets continously
func (coll *Collector) sender() {
	ticker := time.NewTicker(coll.interval)
	for {
		select {
		case <-coll.stop:
			ticker.Stop()
			return
		case <-ticker.C:
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
	var nodeID string
	if val := res.NodeInfo; val != nil {
		nodeID = val.NodeID
	} else if val := res.Neighbours; val != nil {
		nodeID = val.NodeID
	} else if val := res.Statistics; val != nil {
		nodeID = val.NodeID
	}

	// Check length of nodeID
	if len(nodeID) != 12 {
		log.Printf("invalid NodeID '%s' from %s", nodeID, addr.String())
		return
	}

	// Set fields to nil if nodeID is inconsistent
	if res.Statistics != nil && res.Statistics.NodeID != nodeID {
		res.Statistics = nil
	}
	if res.Neighbours != nil && res.Neighbours.NodeID != nodeID {
		res.Neighbours = nil
	}
	if res.NodeInfo != nil && res.NodeInfo.NodeID != nodeID {
		res.NodeInfo = nil
	}

	// Process the data and update IP address
	node := coll.nodes.Update(nodeID, res)
	node.Address = addr.IP

	// Store statistics in database
	if coll.db != nil {
		coll.db.InsertNode(node)
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
			ticker.Stop()
			return
		case <-ticker.C:
			coll.saveGlobalStats()
		}
	}
}

// saves global statistics
func (coll *Collector) saveGlobalStats() {
	stats := runtime.NewGlobalStats(coll.nodes)

	coll.db.InsertGlobals(stats, time.Now())
}
