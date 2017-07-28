package main

import (
	"bufio"
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/respond"
)

const maxDataGramSize = 8192

// Collector for a specificle respond messages
type Collector struct {
	connection *net.UDPConn           // UDP socket
	queue      chan *respond.Response // received responses
	interval   time.Duration          // Interval for multicast packets
	stop       chan interface{}
	nodes      map[string]*data.ResponseData
	interMac   map[string]string
	addrFrom   net.UDPAddr
	addrTo     net.UDPAddr
}

func main() {
	iface := os.Args[1]
	addrFrom := os.Args[2]
	addrTo := os.Args[3]
	linkLocalAddr, err := getLinkLocalAddr(iface)
	if err != nil {
		log.Panic(err)
	}
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   linkLocalAddr,
		Zone: iface,
	})
	if err != nil {
		log.Panic(err)
	}
	conn.SetReadBuffer(maxDataGramSize)
	collector := &Collector{
		connection: conn,
		queue:      make(chan *respond.Response, 400),
		stop:       make(chan interface{}),
		addrFrom:   net.UDPAddr{IP: net.ParseIP(addrFrom)},
		addrTo:     net.UDPAddr{IP: net.ParseIP(addrTo)},
		interval:   time.Second * 10,
		nodes:      make(map[string]*data.ResponseData),
		interMac:   make(map[string]string),
	}
	go collector.receiver(conn)
	go collector.parser()
	collector.sendOnce()
	collector.sender()
	collector.Close()
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

// SendPacket sends a UDP request to the given unicast or multicast address
func (coll *Collector) SendRequestPacket(addr net.UDPAddr) {
	addr.Port = 1001
	if _, err := coll.connection.WriteToUDP([]byte("GET nodeinfo statistics neighbours"), &addr); err != nil {
		log.Println("WriteToUDP failed:", err)
	}
}

func (coll *Collector) saveResponse(addr net.UDPAddr, node *data.ResponseData) {
	if val := node.NodeInfo; val == nil {
		log.Printf("no nodeinfo from %s", addr.String())
		return
	}
	// save current node
	coll.nodes[addr.IP.String()] = node

	// Process the data and update IP address
	var otherIP string
	if addr.IP.Equal(coll.addrFrom.IP) {
		otherIP = coll.addrTo.IP.String()
	} else {
		otherIP = coll.addrFrom.IP.String()
	}

	otherNode := coll.nodes[otherIP]
	if otherIP == "" || otherNode == nil {
		log.Print("othernode not found")
		return
	}

	if node.Neighbours == nil {
		node.Neighbours = &data.Neighbours{
			Batadv: make(map[string]data.BatadvNeighbours),
			NodeID: node.NodeInfo.NodeID,
		}
	}
	interMac := node.NodeInfo.Network.Mesh["bat0"].Interfaces.Other[0]
	if newMac, ok := coll.interMac[addr.IP.String()]; ok {
		interMac = newMac
	} else {
		coll.interMac[addr.IP.String()] = interMac
	}
	if _, ok := node.Neighbours.Batadv[interMac]; !ok {
		node.Neighbours.Batadv[interMac] = data.BatadvNeighbours{
			Neighbours: make(map[string]data.BatmanLink),
		}
	}
	interOtherMac := otherNode.NodeInfo.Network.Mesh["bat0"].Interfaces.Other[0]
	if newMac, ok := coll.interMac[coll.addrTo.IP.String()]; ok {
		interOtherMac = newMac
	} else {
		coll.interMac[otherIP] = interMac
	}
	node.Neighbours.Batadv[interMac].Neighbours[interOtherMac] = data.BatmanLink{
		Tq:       253,
		Lastseen: 0.2,
	}
	buf := bytes.Buffer{}
	writer := bufio.NewWriter(&buf)
	deflater, err := flate.NewWriter(writer, flate.DefaultCompression)

	err = json.NewEncoder(deflater).Encode(node)
	if err != nil {
		panic(err)
	}
	deflater.Close()
	writer.Flush()

	coll.connection.WriteToUDP(buf.Bytes(), &net.UDPAddr{
		IP:   net.ParseIP("fe80::de:faff:fe9f:2414"),
		Port: 12345,
	})
	log.Print("send response from: ", addr.IP.String())
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

		coll.queue <- &respond.Response{
			Address: *src,
			Raw:     raw,
		}
	}
}

func (coll *Collector) parser() {
	for obj := range coll.queue {
		if data, err := obj.Parse(); err != nil {
			log.Println("unable to decode response from", obj.Address.String(), err, "\n", string(obj.Raw))
		} else {
			coll.saveResponse(obj.Address, data)
		}
	}
}

func (coll *Collector) sendOnce() {
	coll.SendRequestPacket(coll.addrFrom)
	coll.SendRequestPacket(coll.addrTo)
	log.Print("send request")
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

// Close Collector
func (coll *Collector) Close() {
	close(coll.stop)
	coll.connection.Close()
	close(coll.queue)
}
