package webserver

import (
	"fmt"
	"strconv"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	VERSION    = ""
	promDescUP = prometheus.NewDesc("yanic_up", "yanic is up", []string{"version"}, prometheus.Labels{})

	promDescNodeInfo = prometheus.NewDesc("yanic_node_info", "node info on labels", []string{
		"node_id",
		"hostname",
		"owner",
		"latitude",
		"longitude",
		"site_code",
		"domain_code",
		"model",
		"nproc",
		"firmware_base",
		"firmware_release",
		"autoupdater",
	}, prometheus.Labels{})

	nodeLabels = []string{"node_id"}

	promDescNodeUP        = prometheus.NewDesc("yanic_node_up", "node is up", nodeLabels, prometheus.Labels{})
	promDescNodeFirstseen = prometheus.NewDesc("yanic_node_firstseen", "node firstseen as timestemp", nodeLabels, prometheus.Labels{})
	promDescNodeLastseen  = prometheus.NewDesc("yanic_node_lastseen", "node lastseen as timestemp", nodeLabels, prometheus.Labels{})

	promDescNodeLoad = prometheus.NewDesc("yanic_node_load", "load of node", nodeLabels, prometheus.Labels{})

	promDescNodeTimeUP   = prometheus.NewDesc("yanic_node_time_up", "", nodeLabels, prometheus.Labels{})
	promDescNodeTimeIdle = prometheus.NewDesc("yanic_node_time_idle", "", nodeLabels, prometheus.Labels{})

	promDescNodeProcRunning = prometheus.NewDesc("yanic_node_proc_running", "", nodeLabels, prometheus.Labels{})

	// Clients - complete
	promDescNodeClientsWifi   = prometheus.NewDesc("yanic_node_clients_wifi", "", nodeLabels, prometheus.Labels{})
	promDescNodeClientsWifi24 = prometheus.NewDesc("yanic_node_clients_wifi24", "", nodeLabels, prometheus.Labels{})
	promDescNodeClientsWifi5  = prometheus.NewDesc("yanic_node_clients_wifi5", "", nodeLabels, prometheus.Labels{})

	promDescNodeClientsOwe   = prometheus.NewDesc("yanic_node_clients_owe", "", nodeLabels, prometheus.Labels{})
	promDescNodeClientsOwe24 = prometheus.NewDesc("yanic_node_clients_owe24", "", nodeLabels, prometheus.Labels{})
	promDescNodeClientsOwe5  = prometheus.NewDesc("yanic_node_clients_owe5", "", nodeLabels, prometheus.Labels{})
	promDescNodeClientsTotal = prometheus.NewDesc("yanic_node_clients_total", "", nodeLabels, prometheus.Labels{})

	// Memory - compelte
	promDescNodeMemBuffers   = prometheus.NewDesc("yanic_node_memmory_buffers", "", nodeLabels, prometheus.Labels{})
	promDescNodeMemCached    = prometheus.NewDesc("yanic_node_memmory_cached", "", nodeLabels, prometheus.Labels{})
	promDescNodeMemFree      = prometheus.NewDesc("yanic_node_memmory_free", "", nodeLabels, prometheus.Labels{})
	promDescNodeMemTotal     = prometheus.NewDesc("yanic_node_memmory_total", "", nodeLabels, prometheus.Labels{})
	promDescNodeMemAvailable = prometheus.NewDesc("yanic_node_memmory_available", "", nodeLabels, prometheus.Labels{})

	// ProcStats - complete
	promDescNodeStatCPUUser    = prometheus.NewDesc("yanic_node_stat_cpu_user", "", nodeLabels, prometheus.Labels{})
	promDescNodeStatCPUNice    = prometheus.NewDesc("yanic_node_stat_cpu_nice", "", nodeLabels, prometheus.Labels{})
	promDescNodeStatCPUSystem  = prometheus.NewDesc("yanic_node_stat_cpu_system", "", nodeLabels, prometheus.Labels{})
	promDescNodeStatCPUIdle    = prometheus.NewDesc("yanic_node_stat_cpu_idle", "", nodeLabels, prometheus.Labels{})
	promDescNodeStatCPUIOWait  = prometheus.NewDesc("yanic_node_stat_cpu_iowait", "", nodeLabels, prometheus.Labels{})
	promDescNodeStatCPUIRQ     = prometheus.NewDesc("yanic_node_stat_cpu_irq", "", nodeLabels, prometheus.Labels{})
	promDescNodeStatCPUSoftIRQ = prometheus.NewDesc("yanic_node_stat_cpu_softirq", "", nodeLabels, prometheus.Labels{})

	promDescNodeStatIntr            = prometheus.NewDesc("yanic_node_stat_intr", "", nodeLabels, prometheus.Labels{})
	promDescNodeStatContextSwitches = prometheus.NewDesc("yanic_node_stat_ctxt", "", nodeLabels, prometheus.Labels{})
	promDescNodeStatSoftIRQ         = prometheus.NewDesc("yanic_node_stat_softirq", "", nodeLabels, prometheus.Labels{})
	promDescNodeStatProcesses       = prometheus.NewDesc("yanic_node_stat_processes", "", nodeLabels, prometheus.Labels{})

	// Traffic - complete
	promDescNodeTrafficRxBytes   = prometheus.NewDesc("yanic_node_traffic_rx_bytes", "", nodeLabels, prometheus.Labels{})
	promDescNodeTrafficRxPackets = prometheus.NewDesc("yanic_node_traffic_rx_packets", "", nodeLabels, prometheus.Labels{})
	promDescNodeTrafficRxDropped = prometheus.NewDesc("yanic_node_traffic_rx_dropped", "", nodeLabels, prometheus.Labels{})

	promDescNodeTrafficTxBytes   = prometheus.NewDesc("yanic_node_traffic_tx_bytes", "", nodeLabels, prometheus.Labels{})
	promDescNodeTrafficTxPackets = prometheus.NewDesc("yanic_node_traffic_tx_packets", "", nodeLabels, prometheus.Labels{})
	promDescNodeTrafficTxDropped = prometheus.NewDesc("yanic_node_traffic_tx_dropped", "", nodeLabels, prometheus.Labels{})

	promDescNodeTrafficForwardBytes   = prometheus.NewDesc("yanic_node_traffic_forward_bytes", "", nodeLabels, prometheus.Labels{})
	promDescNodeTrafficForwardPackets = prometheus.NewDesc("yanic_node_traffic_forward_packets", "", nodeLabels, prometheus.Labels{})
	promDescNodeTrafficForwardDropped = prometheus.NewDesc("yanic_node_traffic_forward_dropped", "", nodeLabels, prometheus.Labels{})

	promDescNodeTrafficMgmtRxBytes   = prometheus.NewDesc("yanic_node_traffic_mgmt_rx_bytes", "", nodeLabels, prometheus.Labels{})
	promDescNodeTrafficMgmtRxPackets = prometheus.NewDesc("yanic_node_traffic_mgmt_rx_packets", "", nodeLabels, prometheus.Labels{})
	promDescNodeTrafficMgmtRxDropped = prometheus.NewDesc("yanic_node_traffic_mgmt_rx_dropped", "", nodeLabels, prometheus.Labels{})

	promDescNodeTrafficMgmtTxBytes   = prometheus.NewDesc("yanic_node_traffic_mgmt_tx_bytes", "", nodeLabels, prometheus.Labels{})
	promDescNodeTrafficMgmtTxPackets = prometheus.NewDesc("yanic_node_traffic_mgmt_tx_packets", "", nodeLabels, prometheus.Labels{})
	promDescNodeTrafficMgmtTxDropped = prometheus.NewDesc("yanic_node_traffic_mgmt_tx_dropped", "", nodeLabels, prometheus.Labels{})

	// Wireless - just necessary
	allLabelsWithFrequency      = append(nodeLabels, "frequency", "frequency_name")
	promDescNodeFrequency       = prometheus.NewDesc("yanic_node_frequency", "", allLabelsWithFrequency, prometheus.Labels{})
	promDescNodeAirtimeChanUtil = prometheus.NewDesc("yanic_node_airtime_chan_util", "", allLabelsWithFrequency, prometheus.Labels{})
	promDescNodeAirtimeTxUtil   = prometheus.NewDesc("yanic_node_airtime_tx_util", "", allLabelsWithFrequency, prometheus.Labels{})
	promDescNodeAirtimeRxUtil   = prometheus.NewDesc("yanic_node_airtime_rx_util", "", allLabelsWithFrequency, prometheus.Labels{})
	promDescNodeAirtimeNoise    = prometheus.NewDesc("yanic_node_airtime_noise", "", allLabelsWithFrequency, prometheus.Labels{})
	promDescNodeWirelessTxPower = prometheus.NewDesc("yanic_node_wireless_txpower", "", allLabelsWithFrequency, prometheus.Labels{})

	// Links - just necessary
	labelLinks     = []string{"source_id", "target_id", "source_address", "target_address", "type"}
	promDescLinkTQ = prometheus.NewDesc("yanic_link_tq", "", labelLinks, prometheus.Labels{})
)

type Prometheus struct {
	Enable        bool     `toml:"enable"`
	DisableLabels []string `toml:"disable_labels"`
	nodes         *runtime.Nodes
	disLabels     map[string]bool
}

func (prom *Prometheus) Init(nodes *runtime.Nodes) {
	prom.nodes = nodes
	prom.disLabels = make(map[string]bool)
	for _, label := range prom.DisableLabels {
		prom.disLabels[label] = false
	}

}

func (prom *Prometheus) Describe(d chan<- *prometheus.Desc) {
	d <- promDescUP
	d <- promDescNodeInfo

	d <- promDescNodeUP
	d <- promDescNodeFirstseen
	d <- promDescNodeLastseen

	d <- promDescNodeLoad

	d <- promDescNodeTimeUP
	d <- promDescNodeTimeIdle

	d <- promDescNodeProcRunning

	d <- promDescNodeClientsWifi
	d <- promDescNodeClientsWifi24
	d <- promDescNodeClientsWifi5

	d <- promDescNodeClientsOwe
	d <- promDescNodeClientsOwe24
	d <- promDescNodeClientsOwe5
	d <- promDescNodeClientsTotal

	d <- promDescNodeMemBuffers
	d <- promDescNodeMemCached
	d <- promDescNodeMemFree
	d <- promDescNodeMemTotal
	d <- promDescNodeMemAvailable

	d <- promDescNodeStatCPUUser
	d <- promDescNodeStatCPUNice
	d <- promDescNodeStatCPUSystem
	d <- promDescNodeStatCPUIdle
	d <- promDescNodeStatCPUIOWait
	d <- promDescNodeStatCPUIRQ
	d <- promDescNodeStatCPUSoftIRQ

	d <- promDescNodeStatIntr
	d <- promDescNodeStatContextSwitches
	d <- promDescNodeStatSoftIRQ
	d <- promDescNodeStatProcesses

	d <- promDescNodeTrafficRxBytes
	d <- promDescNodeTrafficRxPackets
	d <- promDescNodeTrafficRxDropped

	d <- promDescNodeTrafficTxBytes
	d <- promDescNodeTrafficTxPackets
	d <- promDescNodeTrafficTxDropped

	d <- promDescNodeTrafficForwardBytes
	d <- promDescNodeTrafficForwardPackets
	d <- promDescNodeTrafficForwardDropped

	d <- promDescNodeTrafficMgmtRxBytes
	d <- promDescNodeTrafficMgmtRxPackets
	d <- promDescNodeTrafficMgmtRxDropped

	d <- promDescNodeTrafficMgmtTxBytes
	d <- promDescNodeTrafficMgmtTxPackets
	d <- promDescNodeTrafficMgmtTxDropped

	d <- promDescNodeFrequency
	d <- promDescNodeAirtimeChanUtil
	d <- promDescNodeAirtimeRxUtil
	d <- promDescNodeAirtimeTxUtil
	d <- promDescNodeAirtimeNoise
	d <- promDescNodeWirelessTxPower

	d <- promDescLinkTQ
}

func (prom *Prometheus) getNodeLabels(ni *data.Nodeinfo) []string {
	labels := make([]string, 12)
	labels[0] = ni.NodeID
	if _, nok := prom.disLabels["hostname"]; !nok {
		labels[1] = ni.Hostname
	} else {
		labels[1] = "DISABLED"
	}
	if _, nok := prom.disLabels["owner"]; !nok {
		if owner := ni.Owner; owner != nil {
			labels[2] = owner.Contact
		}
	} else {
		labels[2] = "DISABLED"
	}
	if _, nok := prom.disLabels["location"]; !nok {
		if location := ni.Location; location != nil {
			labels[3] = fmt.Sprintf("%v", location.Latitude)
			labels[4] = fmt.Sprintf("%v", location.Longitude)
		}
	} else {
		labels[3] = "DISABLED"
		labels[4] = "DISABLED"
	}
	labels[5] = ni.System.SiteCode
	labels[6] = ni.System.DomainCode
	labels[7] = ni.Hardware.Model
	labels[8] = strconv.Itoa(ni.Hardware.Nproc)
	if firmware := ni.Software.Firmware; firmware != nil {
		labels[9] = firmware.Base
		labels[10] = firmware.Release
	}
	if ni.Software.Autoupdater != nil && ni.Software.Autoupdater.Enabled {
		labels[11] = ni.Software.Autoupdater.Branch
	} else {
		labels[11] = runtime.DISABLED_AUTOUPDATER
	}
	return labels

}

func (prom *Prometheus) Collect(metrics chan<- prometheus.Metric) {
	prom.nodes.Lock()
	defer prom.nodes.Unlock()
	if VERSION != "" {
		if m, err := prometheus.NewConstMetric(
			promDescUP,
			prometheus.GaugeValue,
			1,
			VERSION); err == nil {
			metrics <- m
		}
	}
	for _, node := range prom.nodes.List {
		nodeinfo := node.Nodeinfo
		if nodeinfo == nil {
			continue
		}
		nodeLabelValues := prom.getNodeLabels(nodeinfo)
		if m, err := prometheus.NewConstMetric(
			promDescNodeInfo,
			prometheus.CounterValue,
			1,
			nodeLabelValues...); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeFirstseen,
			prometheus.CounterValue,
			float64(node.Firstseen.Unix()),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeLastseen,
			prometheus.CounterValue,
			float64(node.Lastseen.Unix()),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}
		if !node.Online {
			if m, err := prometheus.NewConstMetric(
				promDescNodeUP,
				prometheus.GaugeValue,
				0,
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			continue
		}
		stats := node.Statistics
		if stats == nil {
			continue
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeUP,
			prometheus.GaugeValue,
			1,
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}

		if m, err := prometheus.NewConstMetric(
			promDescNodeLoad,
			prometheus.GaugeValue,
			stats.LoadAverage,
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}

		if m, err := prometheus.NewConstMetric(
			promDescNodeTimeUP,
			prometheus.GaugeValue,
			float64(stats.Uptime),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}

		if m, err := prometheus.NewConstMetric(
			promDescNodeTimeIdle,
			prometheus.GaugeValue,
			stats.Idletime,
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}

		if m, err := prometheus.NewConstMetric(
			promDescNodeProcRunning,
			prometheus.GaugeValue,
			float64(stats.Processes.Running),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}

		// Clients
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsWifi,
			prometheus.GaugeValue,
			float64(stats.Clients.Wifi),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsWifi24,
			prometheus.GaugeValue,
			float64(stats.Clients.Wifi24),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsWifi5,
			prometheus.GaugeValue,
			float64(stats.Clients.Wifi5),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsOwe,
			prometheus.GaugeValue,
			float64(stats.Clients.OWE),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsOwe24,
			prometheus.GaugeValue,
			float64(stats.Clients.OWE24),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsOwe5,
			prometheus.GaugeValue,
			float64(stats.Clients.OWE5),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsTotal,
			prometheus.GaugeValue,
			float64(stats.Clients.Total),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}

		// Memory
		if m, err := prometheus.NewConstMetric(
			promDescNodeMemBuffers,
			prometheus.GaugeValue,
			float64(stats.Memory.Buffers),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeMemCached,
			prometheus.GaugeValue,
			float64(stats.Memory.Cached),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeMemFree,
			prometheus.GaugeValue,
			float64(stats.Memory.Free),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeMemTotal,
			prometheus.GaugeValue,
			float64(stats.Memory.Total),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeMemAvailable,
			prometheus.GaugeValue,
			float64(stats.Memory.Available),
			nodeinfo.NodeID); err == nil {
			metrics <- m
		}

		if procstat := stats.ProcStats; procstat != nil {
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUUser,
				prometheus.GaugeValue,
				float64(procstat.CPU.User),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUNice,
				prometheus.GaugeValue,
				float64(procstat.CPU.Nice),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUSystem,
				prometheus.GaugeValue,
				float64(procstat.CPU.System),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUIdle,
				prometheus.GaugeValue,
				float64(procstat.CPU.Idle),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUIOWait,
				prometheus.GaugeValue,
				float64(procstat.CPU.IOWait),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUIRQ,
				prometheus.GaugeValue,
				float64(procstat.CPU.IRQ),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUSoftIRQ,
				prometheus.GaugeValue,
				float64(procstat.CPU.SoftIRQ),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}

			if m, err := prometheus.NewConstMetric(
				promDescNodeStatIntr,
				prometheus.GaugeValue,
				float64(procstat.Intr),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatContextSwitches,
				prometheus.GaugeValue,
				float64(procstat.ContextSwitches),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatSoftIRQ,
				prometheus.GaugeValue,
				float64(procstat.SoftIRQ),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatProcesses,
				prometheus.GaugeValue,
				float64(procstat.Processes),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
		}
		if t := stats.Traffic.Rx; t != nil {
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficRxBytes,
				prometheus.GaugeValue,
				float64(t.Bytes),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficRxPackets,
				prometheus.GaugeValue,
				float64(t.Packets),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficRxDropped,
				prometheus.GaugeValue,
				float64(t.Dropped),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
		}
		if t := stats.Traffic.Tx; t != nil {
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficTxBytes,
				prometheus.GaugeValue,
				float64(t.Bytes),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficTxPackets,
				prometheus.GaugeValue,
				float64(t.Packets),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficTxDropped,
				prometheus.GaugeValue,
				float64(t.Dropped),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
		}
		if t := stats.Traffic.Forward; t != nil {
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficForwardBytes,
				prometheus.GaugeValue,
				float64(t.Bytes),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficForwardPackets,
				prometheus.GaugeValue,
				float64(t.Packets),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficForwardDropped,
				prometheus.GaugeValue,
				float64(t.Dropped),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
		}
		if t := stats.Traffic.MgmtTx; t != nil {
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficMgmtTxBytes,
				prometheus.GaugeValue,
				float64(t.Bytes),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficMgmtTxPackets,
				prometheus.GaugeValue,
				float64(t.Packets),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficMgmtTxDropped,
				prometheus.GaugeValue,
				float64(t.Dropped),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
		}
		if t := stats.Traffic.MgmtRx; t != nil {
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficMgmtRxBytes,
				prometheus.GaugeValue,
				float64(t.Bytes),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficMgmtRxPackets,
				prometheus.GaugeValue,
				float64(t.Packets),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficMgmtRxDropped,
				prometheus.GaugeValue,
				float64(t.Dropped),
				nodeinfo.NodeID); err == nil {
				metrics <- m
			}
		}

		// add a label for frequency_name
		labelIndex := len(nodeLabels)
		labels := []string{nodeinfo.NodeID, "", ""}
		for _, airtime := range stats.Wireless {
			labels[labelIndex] = strconv.Itoa(int(airtime.Frequency))
			labels[labelIndex+1] = airtime.FrequencyName()

			if m, err := prometheus.NewConstMetric(
				promDescNodeFrequency,
				prometheus.GaugeValue,
				float64(airtime.Frequency),
				labels...); err == nil {
				metrics <- m
			}

			if m, err := prometheus.NewConstMetric(
				promDescNodeAirtimeChanUtil,
				prometheus.GaugeValue,
				float64(airtime.ChanUtil),
				labels...); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeAirtimeTxUtil,
				prometheus.GaugeValue,
				float64(airtime.TxUtil),
				labels...); err == nil {
				metrics <- m
			}

			if m, err := prometheus.NewConstMetric(
				promDescNodeAirtimeRxUtil,
				prometheus.GaugeValue,
				float64(airtime.RxUtil),
				labels...); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeAirtimeNoise,
				prometheus.GaugeValue,
				float64(airtime.Noise),
				labels...); err == nil {
				metrics <- m
			}
			if wireless := nodeinfo.Wireless; wireless != nil {
				if airtime.Frequency < 5000 {
					if m, err := prometheus.NewConstMetric(
						promDescNodeWirelessTxPower,
						prometheus.GaugeValue,
						float64(wireless.TxPower24),
						labels...); err == nil {
						metrics <- m
					}
				} else {
					if m, err := prometheus.NewConstMetric(
						promDescNodeWirelessTxPower,
						prometheus.GaugeValue,
						float64(wireless.TxPower5),
						labels...); err == nil {
						metrics <- m
					}
				}
			}
		}
		for _, link := range prom.nodes.NodeLinks(node) {
			if m, err := prometheus.NewConstMetric(
				promDescLinkTQ,
				prometheus.GaugeValue,
				float64(link.TQ),
				// labels:
				link.SourceID,
				link.TargetID,
				link.SourceAddress,
				link.TargetAddress,
				link.Type.String(),
			); err == nil {
				metrics <- m
			}
		}
	}
}
