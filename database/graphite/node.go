package graphite

import (
	"time"

	"github.com/FreifunkBremen/yanic/runtime"
	"github.com/fgrosse/graphigo"
)

// PruneNode implementation of database
func (c *Connection) PruneNodes(deleteAfter time.Duration) {
	// we can't really delete nodes from graphite remotely :(
}

// InsertNode implementation of database
func (c *Connection) InsertNode(node *runtime.Node) {
	var fields []graphigo.Metric

	stats := node.Statistics

	nodeinfo := node.Nodeinfo

	if nodeinfo == nil {
		return
	}

	node_prefix := MeasurementNode + `.` + stats.NodeID + `.` + replaceInvalidChars(nodeinfo.Hostname)

	addField := func(name string, value interface{}) {
		fields = append(fields, graphigo.Metric{Name: node_prefix + "." + name, Value: value})
	}

	vpnInterfaces := make(map[string]bool)
	for _, mIface := range nodeinfo.Network.Mesh {
		for _, tunnel := range mIface.Interfaces.Tunnel {
			vpnInterfaces[tunnel] = true
		}
	}

	if neighbours := node.Neighbours; neighbours != nil {
		vpn := 0

		// protocol: Batman Advance
		batadv := 0
		for mac, batadvNeighbours := range neighbours.Batadv {
			batadv += len(batadvNeighbours.Neighbours)
			if _, ok := vpnInterfaces[mac]; ok {
				vpn += len(batadvNeighbours.Neighbours)
			}
		}
		addField("neighbours.batadv", batadv)

		// protocol: Babel
		babel := 0
		for _, babelNeighbours := range neighbours.Babel {
			babel += len(babelNeighbours.Neighbours)
			if _, ok := vpnInterfaces[babelNeighbours.LinkLocalAddress]; ok {
				vpn += len(babelNeighbours.Neighbours)
			}
		}
		addField("neighbours.babel", babel)

		// protocol: LLDP
		lldp := 0
		for _, lldpNeighbours := range neighbours.LLDP {
			lldp += len(lldpNeighbours)
		}
		addField("neighbours.lldp", lldp)

		addField("neighbours.vpn", vpn)

		// total is the sum of all protocols
		addField("neighbours.total", batadv+babel+lldp)
	}

	if t := stats.Traffic.Rx; t != nil {
		addField("traffic.rx.bytes", int64(t.Bytes))
		addField("traffic.rx.packets", t.Packets)
	}
	if t := stats.Traffic.Tx; t != nil {
		addField("traffic.tx.bytes", int64(t.Bytes))
		addField("traffic.tx.packets", t.Packets)
		addField("traffic.tx.dropped", t.Dropped)
	}
	if t := stats.Traffic.Forward; t != nil {
		addField("traffic.forward.bytes", int64(t.Bytes))
		addField("traffic.forward.packets", t.Packets)
	}
	if t := stats.Traffic.MgmtRx; t != nil {
		addField("traffic.mgmt_rx.bytes", int64(t.Bytes))
		addField("traffic.mgmt_rx.packets", t.Packets)
	}
	if t := stats.Traffic.MgmtTx; t != nil {
		addField("traffic.mgmt_tx.bytes", int64(t.Bytes))
		addField("traffic.mgmt_tx.packets", t.Packets)
	}

	for _, airtime := range stats.Wireless {
		suffix := airtime.FrequencyName()
		addField("airtime"+suffix+".chan_util", airtime.ChanUtil)
		addField("airtime"+suffix+".rx_util", airtime.RxUtil)
		addField("airtime"+suffix+".tx_util", airtime.TxUtil)
		addField("airtime"+suffix+".noise", airtime.Noise)
		addField("airtime"+suffix+".frequency", airtime.Frequency)
	}

	addField("load", stats.LoadAverage)
	addField("nproc", nodeinfo.Hardware.Nproc)
	addField("time.up", int64(stats.Uptime))
	addField("time.idle", int64(stats.Idletime))
	addField("proc.running", stats.Processes.Running)
	addField("clients.wifi", stats.Clients.Wifi)
	addField("clients.wifi24", stats.Clients.Wifi24)
	addField("clients.wifi5", stats.Clients.Wifi5)
	addField("clients.total", stats.Clients.Total)
	addField("memory.buffers", stats.Memory.Buffers)
	addField("memory.cached", stats.Memory.Cached)
	addField("memory.free", stats.Memory.Free)
	addField("memory.total", stats.Memory.Total)
	addField("memory.available", stats.Memory.Available)

	c.addPoint(fields)
}
