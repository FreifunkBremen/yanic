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
	"github.com/FreifunkBremen/yanic/lib/jsontime"
	"github.com/FreifunkBremen/yanic/runtime"
)

// Collector for a specificle respond messages
type Collector struct {
	connections []*net.UDPConn          // UDP sockets
	ifaceToConn map[string]*net.UDPConn // map from interface name to UDP socket
	port        int

	queue        chan *Response // received responses
	db           database.Connection
	nodes        *runtime.Nodes
	sitesDomains map[string][]string
	interval     time.Duration // Interval for multicast packets
	stop         chan interface{}
}

// NewCollector creates a Collector struct
func NewCollector(db database.Connection, nodes *runtime.Nodes, sitesDomains map[string][]string, ifaces []string, port int) *Collector {

	coll := &Collector{
		db:           db,
		nodes:        nodes,
		sitesDomains: sitesDomains,
		port:         port,
		queue:        make(chan *Response, 400),
		stop:         make(chan interface{}),
		ifaceToConn:  make(map[string]*net.UDPConn),
	}

	for _, iface := range ifaces {
		coll.listenUDP(iface)
	}

	go coll.parser()

	if coll.db != nil {
		go coll.globalStatsWorker()
	}

	return coll
}

func (coll *Collector) listenUDP(iface string) {
	if _, found := coll.ifaceToConn[iface]; found {
		log.Panicf("can not listen twice on %s", iface)
	}
	linkLocalAddr, err := getLinkLocalAddr(iface)
	if err != nil {
		log.Panic(err)
	}

	// Open socket
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   linkLocalAddr,
		Port: coll.port,
		Zone: iface,
	})
	if err != nil {
		log.Panic(err)
	}
	conn.SetReadBuffer(maxDataGramSize)

	coll.ifaceToConn[iface] = conn
	coll.connections = append(coll.connections, conn)

	// Start receiver
	go coll.receiver(conn)
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
	for _, conn := range coll.connections {
		conn.Close()
	}
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
	log.Println("sending multicasts")
	for _, conn := range coll.connections {
		coll.sendPacket(conn, multiCastGroup)
	}
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
		conn := coll.ifaceToConn[node.Address.Zone]
		if conn == nil {
			log.Printf("unable to find connection for %s", node.Address.Zone)
			continue
		}
		coll.sendPacket(conn, node.Address.IP)
		time.Sleep(10 * time.Millisecond)
	}
}

// SendPacket sends a UDP request to the given unicast or multicast address on the first UDP socket
func (coll *Collector) SendPacket(destination net.IP) {
	coll.sendPacket(coll.connections[0], destination)
}

// sendPacket sends a UDP request to the given unicast or multicast address on the given UDP socket
func (coll *Collector) sendPacket(conn *net.UDPConn, destination net.IP) {
	addr := net.UDPAddr{
		IP:   destination,
		Port: port,
		Zone: conn.LocalAddr().(*net.UDPAddr).Zone,
	}

	if _, err := conn.WriteToUDP([]byte("GET nodeinfo statistics neighbours"), &addr); err != nil {
		log.Println("WriteToUDP failed:", err)
	}
}

// send packets continuously
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

func (coll *Collector) saveResponse(addr *net.UDPAddr, res *data.ResponseData) {
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
	node.Address = addr

	// Store statistics in database
	if db := coll.db; db != nil {
		db.InsertNode(node)

		// Store link data
		if neighbours := node.Neighbours; neighbours != nil {
			coll.nodes.RLock()
			for _, link := range coll.nodes.NodeLinks(node) {
				db.InsertLink(&link, node.Lastseen.GetTime())
			}
			coll.nodes.RUnlock()
		}
	}
}

func (coll *Collector) receiver(conn *net.UDPConn) {
	buf := make([]byte, maxDataGramSize)
	for {
		n, src, err := conn.ReadFromUDP(buf)

		if err != nil {
			log.Println("ReadFromUDP failed:", err)
			return
		}

		raw := make([]byte, n)
		copy(raw, buf)

		coll.queue <- &Response{
			Address: src,
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
	stats := runtime.NewGlobalStats(coll.nodes, coll.sitesDomains)

	for site, domains := range stats {
		for domain, stat := range domains {
			coll.db.InsertGlobals(stat, time.Now(), site, domain)
		}
	}
}
