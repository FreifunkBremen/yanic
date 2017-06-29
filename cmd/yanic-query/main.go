package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/FreifunkBremen/yanic/respond"
	"github.com/FreifunkBremen/yanic/runtime"
)

// Usage: yanic-query wlp4s0 "[fe80::eade:27ff:dead:beef%wlp4s0]:1001"
func main() {
	iface := os.Args[1]
	dstAddress := os.Args[2]

	log.Printf("Sending request address=%s iface=%s", dstAddress, iface)

	nodes := runtime.NewNodes(&runtime.Config{})

	collector := respond.NewCollector(nil, nodes, iface, iface, iface, 0)
	collector.SendPacket(net.UDPAddr{
		IP:   net.ParseIP(dstAddress),
		Zone: iface,
	})

	time.Sleep(time.Second)

	for id, data := range nodes.List {
		log.Printf("%s: %+v", id, data)
	}

	collector.Close()
}
