package influxdb

import (
	"fmt"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
	models "github.com/influxdata/influxdb/models"

	"github.com/FreifunkBremen/yanic/runtime"
)

// InsertNode implementation of database
func (conn *Connection) InsertNode(node *runtime.Node) {
	tags, fields := buildNodeStats(node)
	conn.addPoint(MeasurementNode, tags, fields, time.Now())
}

func (conn *Connection) PruneNodes(deleteAfter time.Duration) {
	query := fmt.Sprintf("delete from %s where time < now() - %ds", MeasurementNode, deleteAfter/time.Second)
	conn.client.Query(client.NewQuery(query, conn.config.Database(), "m"))
}

// returns tags and fields for InfluxDB
func buildNodeStats(node *runtime.Node) (tags models.Tags, fields models.Fields) {
	stats := node.Statistics

	tags.SetString("nodeid", stats.NodeID)

	fields = map[string]interface{}{
		"load":      stats.LoadAverage,
		"time.up":   int64(stats.Uptime),
		"time.idle": int64(stats.Idletime),
	}

	if clients := stats.Clients; clients != nil {
		fields["clients.wifi"] = clients.Wifi
		fields["clients.wifi24"] = clients.Wifi24
		fields["clients.wifi5"] = clients.Wifi5
		fields["clients.total"] = clients.Total
	}

	if proc := stats.Processes; proc != nil {
		fields["proc.running"] = stats.Processes.Running
	}

	if mem := stats.Memory; mem != nil {
		fields["memory.buffers"] = mem.Buffers
		fields["memory.cached"] = mem.Cached
		fields["memory.free"] = mem.Free
		fields["memory.total"] = mem.Total
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
		if firmware := nodeinfo.Software.Firmware; firmware != nil {
			tags.SetString("firmware_base", firmware.Base)
			tags.SetString("firmware_release", firmware.Release)
		}

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
	if tr := stats.Traffic; tr != nil {
		if t := tr.Rx; t != nil {
			fields["traffic.rx.bytes"] = int64(t.Bytes)
			fields["traffic.rx.packets"] = t.Packets
		}
		if t := tr.Tx; t != nil {
			fields["traffic.tx.bytes"] = int64(t.Bytes)
			fields["traffic.tx.packets"] = t.Packets
			fields["traffic.tx.dropped"] = t.Dropped
		}
		if t := tr.Forward; t != nil {
			fields["traffic.forward.bytes"] = int64(t.Bytes)
			fields["traffic.forward.packets"] = t.Packets
		}
		if t := tr.MgmtRx; t != nil {
			fields["traffic.mgmt_rx.bytes"] = int64(t.Bytes)
			fields["traffic.mgmt_rx.packets"] = t.Packets
		}
		if t := tr.MgmtTx; t != nil {
			fields["traffic.mgmt_tx.bytes"] = int64(t.Bytes)
			fields["traffic.mgmt_tx.packets"] = t.Packets
		}
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

	return
}
