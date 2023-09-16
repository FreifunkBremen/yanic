package influxdb

import (
	"strconv"
	"time"

	"github.com/FreifunkBremen/yanic/runtime"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

// PruneNodes prunes historical per-node data - not nessasary, juse configurate your influxdb2
func (conn *Connection) PruneNodes(deleteAfter time.Duration) {
}

// InsertNode stores statistics and neighbours in the database
func (conn *Connection) InsertNode(node *runtime.Node) {
	stats := node.Statistics
	time := node.Lastseen.GetTime()

	if stats == nil || stats.NodeID == "" {
		return
	}

	p := influxdb.NewPoint(MeasurementNode,
		conn.config.Tags(),
		map[string]interface{}{
			"load":             stats.LoadAverage,
			"time.up":          int64(stats.Uptime),
			"time.idle":        int64(stats.Idletime),
			"proc.running":     stats.Processes.Running,
			"clients.wifi":     stats.Clients.Wifi,
			"clients.wifi24":   stats.Clients.Wifi24,
			"clients.wifi5":    stats.Clients.Wifi5,
			"clients.owe":      stats.Clients.OWE,
			"clients.owe24":    stats.Clients.OWE24,
			"clients.owe5":     stats.Clients.OWE5,
			"clients.total":    stats.Clients.Total,
			"memory.buffers":   stats.Memory.Buffers,
			"memory.cached":    stats.Memory.Cached,
			"memory.free":      stats.Memory.Free,
			"memory.total":     stats.Memory.Total,
			"memory.available": stats.Memory.Available,
		},
		time).
		AddTag("nodeid", stats.NodeID)

	vpnInterfaces := make(map[string]bool)

	if nodeinfo := node.Nodeinfo; nodeinfo != nil {
		for _, mIface := range nodeinfo.Network.Mesh {
			for _, tunnel := range mIface.Interfaces.Tunnel {
				vpnInterfaces[tunnel] = true
			}
		}

		p.AddTag("hostname", nodeinfo.Hostname)
		if nodeinfo.System.SiteCode != "" {
			p.AddTag("site", nodeinfo.System.SiteCode)
		}
		if nodeinfo.System.DomainCode != "" {
			p.AddTag("domain", nodeinfo.System.DomainCode)
		}
		if owner := nodeinfo.Owner; owner != nil {
			p.AddTag("owner", owner.Contact)
		}
		if wireless := nodeinfo.Wireless; wireless != nil {
			p.AddField("wireless.txpower24", wireless.TxPower24)
			p.AddField("wireless.txpower5", wireless.TxPower5)
		}
		// Hardware
		p.AddTag("model", nodeinfo.Hardware.Model)
		p.AddField("nproc", nodeinfo.Hardware.Nproc)
		if nodeinfo.Software.Firmware != nil {
			p.AddTag("firmware_base", nodeinfo.Software.Firmware.Base)
			p.AddTag("firmware_release", nodeinfo.Software.Firmware.Release)
		}
		if nodeinfo.Software.Autoupdater != nil && nodeinfo.Software.Autoupdater.Enabled {
			p.AddTag("autoupdater", nodeinfo.Software.Autoupdater.Branch)
		} else {
			p.AddTag("autoupdater", runtime.DISABLED_AUTOUPDATER)
		}

	}
	if neighbours := node.Neighbours; neighbours != nil {
		// VPN Neighbours are Neighbours but includet in one protocol
		vpn := 0

		// protocol: Batman Advance
		batadv := 0
		for mac, batadvNeighbours := range neighbours.Batadv {
			batadv += len(batadvNeighbours.Neighbours)
			if _, ok := vpnInterfaces[mac]; ok {
				vpn += len(batadvNeighbours.Neighbours)
			}
		}
		p.AddField("neighbours.batadv", batadv)

		// protocol: Babel
		babel := 0
		for _, babelNeighbours := range neighbours.Babel {
			babel += len(babelNeighbours.Neighbours)
			if _, ok := vpnInterfaces[babelNeighbours.LinkLocalAddress]; ok {
				vpn += len(babelNeighbours.Neighbours)
			}
		}
		p.AddField("neighbours.babel", babel)

		// protocol: LLDP
		lldp := 0
		for _, lldpNeighbours := range neighbours.LLDP {
			lldp += len(lldpNeighbours)
		}
		p.AddField("neighbours.lldp", lldp)

		// vpn  wait for babel
		p.AddField("neighbours.vpn", vpn)

		// total is the sum of all protocols
		p.AddField("neighbours.total", batadv+babel+lldp)
	}
	if procstat := stats.ProcStats; procstat != nil {
		p.AddField("stat.cpu.user", procstat.CPU.User)
		p.AddField("stat.cpu.nice", procstat.CPU.Nice)
		p.AddField("stat.cpu.system", procstat.CPU.System)
		p.AddField("stat.cpu.idle", procstat.CPU.Idle)
		p.AddField("stat.cpu.iowait", procstat.CPU.IOWait)
		p.AddField("stat.cpu.irq", procstat.CPU.IRQ)
		p.AddField("stat.cpu.softirq", procstat.CPU.SoftIRQ)
		p.AddField("stat.intr", procstat.Intr)
		p.AddField("stat.ctxt", procstat.ContextSwitches)
		p.AddField("stat.softirq", procstat.SoftIRQ)
		p.AddField("stat.processes", procstat.Processes)
	}

	if t := stats.Traffic.Rx; t != nil {
		p.AddField("traffic.rx.bytes", int64(t.Bytes))
		p.AddField("traffic.rx.packets", t.Packets)
	}
	if t := stats.Traffic.Tx; t != nil {
		p.AddField("traffic.tx.bytes", int64(t.Bytes))
		p.AddField("traffic.tx.packets", t.Packets)
		p.AddField("traffic.tx.dropped", t.Dropped)
	}
	if t := stats.Traffic.Forward; t != nil {
		p.AddField("traffic.forward.bytes", int64(t.Bytes))
		p.AddField("traffic.forward.packets", t.Packets)
	}
	if t := stats.Traffic.MgmtRx; t != nil {
		p.AddField("traffic.mgmt_rx.bytes", int64(t.Bytes))
		p.AddField("traffic.mgmt_rx.packets", t.Packets)
	}
	if t := stats.Traffic.MgmtTx; t != nil {
		p.AddField("traffic.mgmt_tx.bytes", int64(t.Bytes))
		p.AddField("traffic.mgmt_tx.packets", t.Packets)
	}

	for _, airtime := range stats.Wireless {
		suffix := airtime.FrequencyName()
		p.AddField("airtime"+suffix+".chan_util", airtime.ChanUtil)
		p.AddField("airtime"+suffix+".rx_util", airtime.RxUtil)
		p.AddField("airtime"+suffix+".tx_util", airtime.TxUtil)
		p.AddField("airtime"+suffix+".noise", airtime.Noise)
		p.AddField("airtime"+suffix+".frequency", airtime.Frequency)
		p.AddTag("frequency"+suffix, strconv.Itoa(int(airtime.Frequency)))
	}

	conn.writeAPI[MeasurementNode].WritePoint(p)

	// Add DHCP statistics
	if dhcp := stats.DHCP; dhcp != nil {
		p := influxdb.NewPoint(MeasurementDHCP,
			conn.config.Tags(),
			map[string]interface{}{
				"decline":  dhcp.Decline,
				"offer":    dhcp.Offer,
				"ack":      dhcp.Ack,
				"nak":      dhcp.Nak,
				"request":  dhcp.Request,
				"discover": dhcp.Discover,
				"inform":   dhcp.Inform,
				"release":  dhcp.Release,

				"leases.allocated": dhcp.LeasesAllocated,
				"leases.pruned":    dhcp.LeasesPruned,
			}, time).
			AddTag("nodeid", stats.NodeID)

		if nodeinfo := node.Nodeinfo; nodeinfo != nil {
			p.AddTag("hostname", nodeinfo.Hostname)
		}

		conn.writeAPI[MeasurementDHCP].WritePoint(p)
	}
}
