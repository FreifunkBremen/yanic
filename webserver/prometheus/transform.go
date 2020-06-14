package prometheus

import (
	"github.com/FreifunkBremen/yanic/runtime"
)

func MetricLabelsFromNode(node *runtime.Node) (labels map[string]interface{}) {
	labels = make(map[string]interface{})

	nodeinfo := node.Nodeinfo
	if nodeinfo == nil {
		return
	}

	labels["node_id"] =  nodeinfo.NodeID
	labels["hostname"] = nodeinfo.Hostname

	if nodeinfo.System.SiteCode != "" {
		labels["site_code"] = nodeinfo.System.SiteCode
	}
	if nodeinfo.System.DomainCode != "" {
		labels["domain_code"] = nodeinfo.System.DomainCode
	}
	if owner := nodeinfo.Owner; owner != nil {
		labels["owner"] = owner.Contact
	}
	// Hardware
	labels["model"] = nodeinfo.Hardware.Model
	labels["nproc"] = nodeinfo.Hardware.Nproc
	labels["firmware_base"] = nodeinfo.Software.Firmware.Base
	labels["firmware_release"] = nodeinfo.Software.Firmware.Release
	if nodeinfo.Software.Autoupdater.Enabled {
		labels["autoupdater"] = nodeinfo.Software.Autoupdater.Branch
	} else {
		labels["autoupdater"] = runtime.DISABLED_AUTOUPDATER
	}

	if location := nodeinfo.Location; location != nil {
		labels["location_lat"] = location.Latitude
		labels["location_long"] = location.Longitude
	}

	return
}

func MetricsFromNode(nodes *runtime.Nodes, node *runtime.Node) []Metric {
	m := []Metric{}

	// before node metrics to get link statics undependent of node validation
	for _, link := range nodes.NodeLinks(node) {
		m = append(m, Metric{
			Labels: map[string]interface{}{
				"source_id":   link.SourceID,
				"source_addr": link.SourceAddress,
				"target_id":   link.TargetID,
				"target_addr": link.TargetAddress,
			},
			Name:  "yanic_link",
			Value: link.TQ * 100,
		})
	}

	nodeinfo := node.Nodeinfo
	stats := node.Statistics

	// validation
	if nodeinfo == nil || stats == nil {
		return m
	}

	labels := MetricLabelsFromNode(node)

	addMetric := func(name string, value interface{}) {
		m = append(m, Metric{Labels: labels, Name: "yanic_" + name, Value: value})
	}

	if node.Online {
		addMetric("node_up", 1)
	} else {
		addMetric("node_up", 0)
	}

	addMetric("node_load", stats.LoadAverage)

	addMetric("node_time_up", stats.Uptime)
	addMetric("node_time_idle", stats.Idletime)

	addMetric("node_proc_running", stats.Processes.Running)

	addMetric("node_clients_wifi", stats.Clients.Wifi)
	addMetric("node_clients_wifi24", stats.Clients.Wifi24)
	addMetric("node_clients_wifi5", stats.Clients.Wifi5)
	addMetric("node_clients_total", stats.Clients.Total)

	addMetric("node_memory_buffers", stats.Memory.Buffers)
	addMetric("node_memory_cached", stats.Memory.Cached)
	addMetric("node_memory_free", stats.Memory.Free)
	addMetric("node_memory_total", stats.Memory.Total)
	addMetric("node_memory_available", stats.Memory.Available)

	//TODO Neighbours count after merging improvement in influxdb and graphite

	if procstat := stats.ProcStats; procstat != nil {
		addMetric("node_stat_cpu_user", procstat.CPU.User)
		addMetric("node_stat_cpu_nice", procstat.CPU.Nice)
		addMetric("node_stat_cpu_system", procstat.CPU.System)
		addMetric("node_stat_cpu_idle", procstat.CPU.Idle)
		addMetric("node_stat_cpu_iowait", procstat.CPU.IOWait)
		addMetric("node_stat_cpu_irq", procstat.CPU.IRQ)
		addMetric("node_stat_cpu_softirq", procstat.CPU.SoftIRQ)
		addMetric("node_stat_intr", procstat.Intr)
		addMetric("node_stat_ctxt", procstat.ContextSwitches)
		addMetric("node_stat_softirq", procstat.SoftIRQ)
		addMetric("node_stat_processes", procstat.Processes)
	}

	if t := stats.Traffic.Rx; t != nil {
		addMetric("node_traffic_rx_bytes", t.Bytes)
		addMetric("node_traffic_rx_packets", t.Packets)
	}
	if t := stats.Traffic.Tx; t != nil {
		addMetric("node_traffic_tx_bytes", t.Bytes)
		addMetric("node_traffic_tx_packets", t.Packets)
		addMetric("node_traffic_tx_dropped", t.Dropped)
	}
	if t := stats.Traffic.Forward; t != nil {
		addMetric("node_traffic_forward_bytes", t.Bytes)
		addMetric("node_traffic_forward_packets", t.Packets)
	}
	if t := stats.Traffic.MgmtRx; t != nil {
		addMetric("node_traffic_mgmt_rx_bytes", t.Bytes)
		addMetric("node_traffic_mgmt_rx_packets", t.Packets)
	}
	if t := stats.Traffic.MgmtTx; t != nil {
		addMetric("node_traffic_mgmt_tx_bytes", t.Bytes)
		addMetric("node_traffic_mgmt_tx_packets", t.Packets)
	}

	for _, airtime := range stats.Wireless {
		labels["frequency_name"] = airtime.FrequencyName()
		addMetric("node_frequency", airtime.Frequency)
		addMetric("node_airtime_chan_util", airtime.ChanUtil)
		addMetric("node_airtime_rx_util", airtime.RxUtil)
		addMetric("node_airtime_tx_util", airtime.TxUtil)
		addMetric("node_airtime_noise", airtime.Noise)
		if wireless := nodeinfo.Wireless; wireless != nil {
			if airtime.Frequency < 5000 {
				addMetric("node_wireless_txpower", wireless.TxPower24)
			} else {
				addMetric("node_wireless_txpower", wireless.TxPower5)
			}
		}
	}

	return m
}
