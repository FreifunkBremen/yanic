package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/FreifunkBremen/yanic/models"
	"github.com/FreifunkBremen/yanic/respond"
)

// Usage: respond-query wlp4s0 "[fe80::eade:27ff:dead:beef%wlp4s0]:1001"
func main() {
	iface := os.Args[1]
	dstAddress := os.Args[2]

	log.Printf("Sending request address=%s iface=%s", dstAddress, iface)

	nodes := models.NewNodes(&models.Config{})

	collector := respond.NewCollector(nil, nodes, iface, 0)
	collector.SendPacket(net.ParseIP(dstAddress))

	time.Sleep(time.Second)

	for id, data := range nodes.List {
		log.Printf("%s: %+v", id, data)
	}

	collector.Close()
}
