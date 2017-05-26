package graphite

import (
	"strings"
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

	node_prefix := MeasurementNode + `.` + stats.NodeID + `.` + strings.Replace(nodeinfo.Hostname, ".", "__", -1)

	if neighbours := node.Neighbours; neighbours != nil {
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
		fields = append(fields, graphigo.Metric{Name: node_prefix + ".neighbours.vpn", Value: vpn})
		// protocol: Batman Advance
		batadv := 0
		for _, batadvNeighbours := range neighbours.Batadv {
			batadv += len(batadvNeighbours.Neighbours)
		}
		fields = append(fields, graphigo.Metric{Name: node_prefix + ".neighbours.batadv", Value: batadv})

		// protocol: LLDP
		lldp := 0
		for _, lldpNeighbours := range neighbours.LLDP {
			lldp += len(lldpNeighbours)
		}
		fields = append(fields, graphigo.Metric{Name: node_prefix + ".neighbours.lldp", Value: lldp})

		// total is the sum of all protocols
		fields = append(fields, graphigo.Metric{Name: node_prefix + ".neighbours.total", Value: batadv + lldp})
	}

	if t := stats.Traffic.Rx; t != nil {
		fields = append(fields, graphigo.Metric{Name: node_prefix + ".traffic.rx.bytes", Value: int64(t.Bytes)},
			graphigo.Metric{Name: node_prefix + ".traffic.rx.packets", Value: t.Packets})
	}
	if t := stats.Traffic.Tx; t != nil {
		fields = append(fields, graphigo.Metric{Name: node_prefix + ".traffic.tx.bytes", Value: int64(t.Bytes)},
			graphigo.Metric{Name: node_prefix + ".traffic.tx.packets", Value: t.Packets},
			graphigo.Metric{Name: node_prefix + ".traffic.tx.dropped", Value: t.Dropped})
	}
	if t := stats.Traffic.Forward; t != nil {
		fields = append(fields, graphigo.Metric{Name: node_prefix + ".traffic.forward.bytes", Value: int64(t.Bytes)},
			graphigo.Metric{Name: node_prefix + ".traffic.forward.packets", Value: t.Packets})
	}
	if t := stats.Traffic.MgmtRx; t != nil {
		fields = append(fields, graphigo.Metric{Name: node_prefix + ".traffic.mgmt_rx.bytes", Value: int64(t.Bytes)},
			graphigo.Metric{Name: node_prefix + ".traffic.mgmt_rx.packets", Value: t.Packets})
	}
	if t := stats.Traffic.MgmtTx; t != nil {
		fields = append(fields, graphigo.Metric{Name: node_prefix + ".traffic.mgmt_tx.bytes", Value: int64(t.Bytes)},
			graphigo.Metric{Name: node_prefix + ".traffic.mgmt_tx.packets", Value: t.Packets})
	}

	for _, airtime := range stats.Wireless {
		suffix := airtime.FrequencyName()
		fields = append(fields, []graphigo.Metric{
			{Name: node_prefix + ".airtime" + suffix + ".chan_util", Value: airtime.ChanUtil},
			{Name: node_prefix + ".airtime" + suffix + ".rx_util", Value: airtime.RxUtil},
			{Name: node_prefix + ".airtime" + suffix + ".tx_util", Value: airtime.TxUtil},
			{Name: node_prefix + ".airtime" + suffix + ".noise", Value: airtime.Noise},
			{Name: node_prefix + ".airtime" + suffix + ".frequency", Value: airtime.Frequency},
		}...)
	}

	fields = append(fields, []graphigo.Metric{
		{Name: node_prefix + ".load", Value: stats.LoadAverage},
		{Name: node_prefix + ".time.up", Value: int64(stats.Uptime)},
		{Name: node_prefix + ".time.idle", Value: int64(stats.Idletime)},
		{Name: node_prefix + ".proc.running", Value: stats.Processes.Running},
		{Name: node_prefix + ".clients.wifi", Value: stats.Clients.Wifi},
		{Name: node_prefix + ".clients.wifi24", Value: stats.Clients.Wifi24},
		{Name: node_prefix + ".clients.wifi5", Value: stats.Clients.Wifi5},
		{Name: node_prefix + ".clients.total", Value: stats.Clients.Total},
		{Name: node_prefix + ".memory.buffers", Value: stats.Memory.Buffers},
		{Name: node_prefix + ".memory.cached", Value: stats.Memory.Cached},
		{Name: node_prefix + ".memory.free", Value: stats.Memory.Free},
		{Name: node_prefix + ".memory.total", Value: stats.Memory.Total},
	}...)

	c.addPoint(fields)
}
