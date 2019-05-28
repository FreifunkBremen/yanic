package influxdb

import (
	"fmt"
	"strconv"
	"time"

	models "github.com/influxdata/influxdb1-client/models"
	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/FreifunkBremen/yanic/runtime"
)

// PruneNodes prunes historical per-node data
func (conn *Connection) PruneNodes(deleteAfter time.Duration) {
	for _, measurement := range []string{MeasurementNode, MeasurementLink} {
		query := fmt.Sprintf("delete from %s where time < now() - %ds", measurement, deleteAfter/time.Second)
		conn.client.Query(client.NewQuery(query, conn.config.Database(), "m"))
	}

}

// InsertNode stores statistics and neighbours in the database
func (conn *Connection) InsertNode(node *runtime.Node) {
	stats := node.Statistics
	time := node.Lastseen.GetTime()

	if stats == nil || stats.NodeID == "" {
		return
	}

	tags := models.Tags{}
	tags.SetString("nodeid", stats.NodeID)

	fields := models.Fields{
		"load":             stats.LoadAverage,
		"time.up":          int64(stats.Uptime),
		"time.idle":        int64(stats.Idletime),
		"proc.running":     stats.Processes.Running,
		"clients.wifi":     stats.Clients.Wifi,
		"clients.wifi24":   stats.Clients.Wifi24,
		"clients.wifi5":    stats.Clients.Wifi5,
		"clients.total":    stats.Clients.Total,
		"memory.buffers":   stats.Memory.Buffers,
		"memory.cached":    stats.Memory.Cached,
		"memory.free":      stats.Memory.Free,
		"memory.total":     stats.Memory.Total,
		"memory.available": stats.Memory.Available,
	}

	vpnInterfaces := make(map[string]bool)

	if nodeinfo := node.Nodeinfo; nodeinfo != nil {
		for _, mIface := range nodeinfo.Network.Mesh {
			for _, tunnel := range mIface.Interfaces.Tunnel {
				vpnInterfaces[tunnel] = true
			}
		}

		tags.SetString("hostname", nodeinfo.Hostname)
		if nodeinfo.System.SiteCode != "" {
			tags.SetString("site", nodeinfo.System.SiteCode)
		}
		if nodeinfo.System.DomainCode != "" {
			tags.SetString("domain", nodeinfo.System.DomainCode)
		}
		if owner := nodeinfo.Owner; owner != nil {
			tags.SetString("owner", owner.Contact)
		}
		if wireless := nodeinfo.Wireless; wireless != nil {
			fields["wireless.txpower24"] = wireless.TxPower24
			fields["wireless.txpower5"] = wireless.TxPower5
		}
		// Hardware
		tags.SetString("model", nodeinfo.Hardware.Model)
		fields["nproc"] = nodeinfo.Hardware.Nproc
		tags.SetString("firmware_base", nodeinfo.Software.Firmware.Base)
		tags.SetString("firmware_release", nodeinfo.Software.Firmware.Release)
		if nodeinfo.Software.Autoupdater.Enabled {
			tags.SetString("autoupdater", nodeinfo.Software.Autoupdater.Branch)
		} else {
			tags.SetString("autoupdater", runtime.DISABLED_AUTOUPDATER)
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
		fields["neighbours.batadv"] = batadv

		// protocol: Babel
		babel := 0
		for _, babelNeighbours := range neighbours.Babel {
			babel += len(babelNeighbours.Neighbours)
			if _, ok := vpnInterfaces[babelNeighbours.LinkLocalAddress]; ok {
				vpn += len(babelNeighbours.Neighbours)
			}
		}
		fields["neighbours.babel"] = babel

		// protocol: LLDP
		lldp := 0
		for _, lldpNeighbours := range neighbours.LLDP {
			lldp += len(lldpNeighbours)
		}
		fields["neighbours.lldp"] = lldp

		// vpn  wait for babel
		fields["neighbours.vpn"] = vpn

		// total is the sum of all protocols
		fields["neighbours.total"] = batadv + babel + lldp
	}
	if procstat := stats.ProcStats; procstat != nil {
		fields["stat.cpu.user"] = procstat.CPU.User
		fields["stat.cpu.nice"] = procstat.CPU.Nice
		fields["stat.cpu.system"] = procstat.CPU.System
		fields["stat.cpu.idle"] = procstat.CPU.Idle
		fields["stat.cpu.iowait"] = procstat.CPU.IOWait
		fields["stat.cpu.irq"] = procstat.CPU.IRQ
		fields["stat.cpu.softirq"] = procstat.CPU.SoftIRQ
		fields["stat.intr"] = procstat.Intr
		fields["stat.ctxt"] = procstat.ContextSwitches
		fields["stat.softirq"] = procstat.SoftIRQ
		fields["stat.processes"] = procstat.Processes
	}

	if t := stats.Traffic.Rx; t != nil {
		fields["traffic.rx.bytes"] = int64(t.Bytes)
		fields["traffic.rx.packets"] = t.Packets
	}
	if t := stats.Traffic.Tx; t != nil {
		fields["traffic.tx.bytes"] = int64(t.Bytes)
		fields["traffic.tx.packets"] = t.Packets
		fields["traffic.tx.dropped"] = t.Dropped
	}
	if t := stats.Traffic.Forward; t != nil {
		fields["traffic.forward.bytes"] = int64(t.Bytes)
		fields["traffic.forward.packets"] = t.Packets
	}
	if t := stats.Traffic.MgmtRx; t != nil {
		fields["traffic.mgmt_rx.bytes"] = int64(t.Bytes)
		fields["traffic.mgmt_rx.packets"] = t.Packets
	}
	if t := stats.Traffic.MgmtTx; t != nil {
		fields["traffic.mgmt_tx.bytes"] = int64(t.Bytes)
		fields["traffic.mgmt_tx.packets"] = t.Packets
	}

	for _, airtime := range stats.Wireless {
		suffix := airtime.FrequencyName()
		fields["airtime"+suffix+".chan_util"] = airtime.ChanUtil
		fields["airtime"+suffix+".rx_util"] = airtime.RxUtil
		fields["airtime"+suffix+".tx_util"] = airtime.TxUtil
		fields["airtime"+suffix+".noise"] = airtime.Noise
		fields["airtime"+suffix+".frequency"] = airtime.Frequency
		tags.SetString("frequency"+suffix, strconv.Itoa(int(airtime.Frequency)))
	}

	conn.addPoint(MeasurementNode, tags, fields, time)

	// Add DHCP statistics
	if dhcp := stats.DHCP; dhcp != nil {
		fields := models.Fields{
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
		}

		// Tags
		tags.SetString("nodeid", stats.NodeID)
		if nodeinfo := node.Nodeinfo; nodeinfo != nil {
			tags.SetString("hostname", nodeinfo.Hostname)
		}

		conn.addPoint(MeasurementDHCP, tags, fields, time)
	}

	return
}
