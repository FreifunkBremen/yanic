package influxdb

import (
	"fmt"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
	models "github.com/influxdata/influxdb/models"

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

	if stats == nil {
		return
	}

	tags := models.Tags{}
	tags.SetString("nodeid", stats.NodeID)

	fields := models.Fields{
		"load":           stats.LoadAverage,
		"time.up":        int64(stats.Uptime),
		"time.idle":      int64(stats.Idletime),
		"proc.running":   stats.Processes.Running,
		"clients.wifi":   stats.Clients.Wifi,
		"clients.wifi24": stats.Clients.Wifi24,
		"clients.wifi5":  stats.Clients.Wifi5,
		"clients.total":  stats.Clients.Total,
		"memory.buffers": stats.Memory.Buffers,
		"memory.cached":  stats.Memory.Cached,
		"memory.free":    stats.Memory.Free,
		"memory.total":   stats.Memory.Total,
	}

	if nodeinfo := node.Nodeinfo; nodeinfo != nil {
		tags.SetString("hostname", nodeinfo.Hostname)
		if owner := nodeinfo.Owner; owner != nil {
			tags.SetString("owner", owner.Contact)
		}
		if wireless := nodeinfo.Wireless; wireless != nil {
			fields["wireless.txpower24"] = wireless.TxPower24
			fields["wireless.txpower5"] = wireless.TxPower5
		}
		// Hardware
		tags.SetString("model", nodeinfo.Hardware.Model)
		tags.SetString("firmware_base", nodeinfo.Software.Firmware.Base)
		tags.SetString("firmware_release", nodeinfo.Software.Firmware.Release)

	}

	if neighbours := node.Neighbours; neighbours != nil {
		// VPN Neighbours are Neighbours but includet in one protocol
		vpn := 0
		if meshvpn := stats.MeshVPN; meshvpn != nil {
			for _, group := range meshvpn.Groups {
				for _, link := range group.Peers {
					if link != nil && link.Established > 1 {
						vpn++
					}
				}
			}
		}
		fields["neighbours.vpn"] = vpn

		// protocol: Batman Advance
		batadv := 0
		for _, batadvNeighbours := range neighbours.Batadv {
			batadv += len(batadvNeighbours.Neighbours)
			for neighbourID, link := range batadvNeighbours.Neighbours {
				conn.insertLinkStatistics(stats.NodeID, neighbourID, link.Tq, time)
			}
		}
		fields["neighbours.batadv"] = batadv

		// protocol: LLDP
		lldp := 0
		for _, lldpNeighbours := range neighbours.LLDP {
			lldp += len(lldpNeighbours)
		}
		fields["neighbours.lldp"] = lldp

		// total is the sum of all protocols
		fields["neighbours.total"] = batadv + lldp
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

	return
}

// adds a link data point
func (conn *Connection) insertLinkStatistics(source string, target string, tq int, t time.Time) {
	tags := models.Tags{}
	tags.SetString("source", source)
	tags.SetString("target", target)

	conn.addPoint(MeasurementLink, tags, models.Fields{"tq": tq}, t)
}
