package webserver

import (
	"fmt"
	"strconv"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	allLabels = []string{"node_id", "hostname", "owner", "latitude", "longitude", "site_code", "domain_code", "model", "nproc", "firmware_base", "firmware_release", "autoupdater"}

	promDescNodeUP = prometheus.NewDesc("yanic_node_up", "is the node there", allLabels, prometheus.Labels{})

	promDescNodeLoad = prometheus.NewDesc("yanic_node_load", "", allLabels, prometheus.Labels{})

	promDescNodeTimeUP   = prometheus.NewDesc("yanic_node_time_up", "", allLabels, prometheus.Labels{})
	promDescNodeTimeIdle = prometheus.NewDesc("yanic_node_time_idle", "", allLabels, prometheus.Labels{})

	promDescNodeProcRunning = prometheus.NewDesc("yanic_node_proc_running", "", allLabels, prometheus.Labels{})

	// Clients - complete
	promDescNodeClientsWifi   = prometheus.NewDesc("yanic_node_clients_wifi", "", allLabels, prometheus.Labels{})
	promDescNodeClientsWifi24 = prometheus.NewDesc("yanic_node_clients_wifi24", "", allLabels, prometheus.Labels{})
	promDescNodeClientsWifi5  = prometheus.NewDesc("yanic_node_clients_wifi5", "", allLabels, prometheus.Labels{})

	promDescNodeClientsOwe   = prometheus.NewDesc("yanic_node_clients_owe", "", allLabels, prometheus.Labels{})
	promDescNodeClientsOwe24 = prometheus.NewDesc("yanic_node_clients_owe24", "", allLabels, prometheus.Labels{})
	promDescNodeClientsOwe5  = prometheus.NewDesc("yanic_node_clients_owe5", "", allLabels, prometheus.Labels{})
	promDescNodeClientsTotal  = prometheus.NewDesc("yanic_node_clients_total", "", allLabels, prometheus.Labels{})

	// Memory - compelte
	promDescNodeMemBuffers   = prometheus.NewDesc("yanic_node_memmory_buffers", "", allLabels, prometheus.Labels{})
	promDescNodeMemCached    = prometheus.NewDesc("yanic_node_memmory_cached", "", allLabels, prometheus.Labels{})
	promDescNodeMemFree      = prometheus.NewDesc("yanic_node_memmory_free", "", allLabels, prometheus.Labels{})
	promDescNodeMemTotal     = prometheus.NewDesc("yanic_node_memmory_total", "", allLabels, prometheus.Labels{})
	promDescNodeMemAvailable = prometheus.NewDesc("yanic_node_memmory_available", "", allLabels, prometheus.Labels{})

	// ProcStats - complete
	promDescNodeStatCPUUser = prometheus.NewDesc("yanic_node_stat_cpu_user", "", allLabels, prometheus.Labels{})
	promDescNodeStatCPUNice = prometheus.NewDesc("yanic_node_stat_cpu_nice", "", allLabels, prometheus.Labels{})
	promDescNodeStatCPUSystem = prometheus.NewDesc("yanic_node_stat_cpu_system", "", allLabels, prometheus.Labels{})
	promDescNodeStatCPUIdle = prometheus.NewDesc("yanic_node_stat_cpu_idle", "", allLabels, prometheus.Labels{})
	promDescNodeStatCPUIOWait = prometheus.NewDesc("yanic_node_stat_cpu_iowait", "", allLabels, prometheus.Labels{})
	promDescNodeStatCPUIRQ = prometheus.NewDesc("yanic_node_stat_cpu_irq", "", allLabels, prometheus.Labels{})
	promDescNodeStatCPUSoftIRQ = prometheus.NewDesc("yanic_node_stat_cpu_softirq", "", allLabels, prometheus.Labels{})

	promDescNodeStatIntr = prometheus.NewDesc("yanic_node_stat_intr", "", allLabels, prometheus.Labels{})
	promDescNodeStatContextSwitches = prometheus.NewDesc("yanic_node_stat_ctxt", "", allLabels, prometheus.Labels{})
	promDescNodeStatSoftIRQ = prometheus.NewDesc("yanic_node_stat_softirq", "", allLabels, prometheus.Labels{})
	promDescNodeStatProcesses = prometheus.NewDesc("yanic_node_stat_processes", "", allLabels, prometheus.Labels{})

	// Traffic - complete
	promDescNodeTrafficRxBytes = prometheus.NewDesc("yanic_node_traffic_rx_bytes", "", allLabels, prometheus.Labels{})
	promDescNodeTrafficRxPackets = prometheus.NewDesc("yanic_node_traffic_rx_packets", "", allLabels, prometheus.Labels{})
	promDescNodeTrafficRxDropped = prometheus.NewDesc("yanic_node_traffic_rx_dropped", "", allLabels, prometheus.Labels{})

	promDescNodeTrafficTxBytes = prometheus.NewDesc("yanic_node_traffic_tx_bytes", "", allLabels, prometheus.Labels{})
	promDescNodeTrafficTxPackets = prometheus.NewDesc("yanic_node_traffic_tx_packets", "", allLabels, prometheus.Labels{})
	promDescNodeTrafficTxDropped = prometheus.NewDesc("yanic_node_traffic_tx_dropped", "", allLabels, prometheus.Labels{})

	promDescNodeTrafficForwardBytes = prometheus.NewDesc("yanic_node_traffic_forward_bytes", "", allLabels, prometheus.Labels{})
	promDescNodeTrafficForwardPackets = prometheus.NewDesc("yanic_node_traffic_forward_packets", "", allLabels, prometheus.Labels{})
	promDescNodeTrafficForwardDropped = prometheus.NewDesc("yanic_node_traffic_forward_dropped", "", allLabels, prometheus.Labels{})

	promDescNodeTrafficMgmtRxBytes = prometheus.NewDesc("yanic_node_traffic_mgmt_rx_bytes", "", allLabels, prometheus.Labels{})
	promDescNodeTrafficMgmtRxPackets = prometheus.NewDesc("yanic_node_traffic_mgmt_rx_packets", "", allLabels, prometheus.Labels{})
	promDescNodeTrafficMgmtRxDropped = prometheus.NewDesc("yanic_node_traffic_mgmt_rx_dropped", "", allLabels, prometheus.Labels{})

	promDescNodeTrafficMgmtTxBytes = prometheus.NewDesc("yanic_node_traffic_mgmt_tx_bytes", "", allLabels, prometheus.Labels{})
	promDescNodeTrafficMgmtTxPackets = prometheus.NewDesc("yanic_node_traffic_mgmt_tx_packets", "", allLabels, prometheus.Labels{})
	promDescNodeTrafficMgmtTxDropped = prometheus.NewDesc("yanic_node_traffic_mgmt_tx_dropped", "", allLabels, prometheus.Labels{})

	// Wireless - just necessary
	allLabelsWithFrequency      = append(allLabels, "frequency", "frequency_name")
	promDescNodeFrequency       = prometheus.NewDesc("yanic_node_frequency", "", allLabelsWithFrequency, prometheus.Labels{})
	promDescNodeAirtimeChanUtil = prometheus.NewDesc("yanic_node_airtime_chan_util", "", allLabelsWithFrequency, prometheus.Labels{})
	promDescNodeAirtimeTxUtil   = prometheus.NewDesc("yanic_node_airtime_tx_util", "", allLabelsWithFrequency, prometheus.Labels{})
	promDescNodeAirtimeRxUtil   = prometheus.NewDesc("yanic_node_airtime_rx_util", "", allLabelsWithFrequency, prometheus.Labels{})
	promDescNodeAirtimeNoise    = prometheus.NewDesc("yanic_node_airtime_noise", "", allLabelsWithFrequency, prometheus.Labels{})
	promDescNodeWirelessTxPower = prometheus.NewDesc("yanic_node_wireless_txpower", "", allLabelsWithFrequency, prometheus.Labels{})
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
	d <- promDescNodeUP

	d <- promDescNodeLoad

	d <- promDescNodeTimeUP
	d <- promDescNodeTimeIdle

	d <-promDescNodeProcRunning

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
}

func (prom *Prometheus) getNodeLabels(ni *data.Nodeinfo) []string {
	labels := make([]string, len(allLabels))
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
	for _, node := range prom.nodes.List {
		nodeinfo := node.Nodeinfo
		if nodeinfo == nil {
			continue
		}
		labels := prom.getNodeLabels(nodeinfo)
		if !node.Online {
			if m, err := prometheus.NewConstMetric(
				promDescNodeUP,
				prometheus.GaugeValue,
				0,
				labels...); err == nil {
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
			labels...); err == nil {
			metrics <- m
		}

		if m, err := prometheus.NewConstMetric(
			promDescNodeLoad,
			prometheus.GaugeValue,
			stats.LoadAverage,
			labels...); err == nil {
			metrics <- m
		}

		if m, err := prometheus.NewConstMetric(
			promDescNodeTimeUP,
			prometheus.GaugeValue,
			float64(stats.Uptime),
			labels...); err == nil {
			metrics <- m
		}

		if m, err := prometheus.NewConstMetric(
			promDescNodeTimeIdle,
			prometheus.GaugeValue,
			stats.Idletime,
			labels...); err == nil {
			metrics <- m
		}

		if m, err := prometheus.NewConstMetric(
			promDescNodeProcRunning,
			prometheus.GaugeValue,
			float64(stats.Processes.Running),
			labels...); err == nil {
			metrics <- m
		}

		// Clients
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsWifi,
			prometheus.GaugeValue,
			float64(stats.Clients.Wifi),
			labels...); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsWifi24,
			prometheus.GaugeValue,
			float64(stats.Clients.Wifi24),
			labels...); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsWifi5,
			prometheus.GaugeValue,
			float64(stats.Clients.Wifi5),
			labels...); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsOwe,
			prometheus.GaugeValue,
			float64(stats.Clients.Owe),
			labels...); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsOwe24,
			prometheus.GaugeValue,
			float64(stats.Clients.Owe24),
			labels...); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsOwe5,
			prometheus.GaugeValue,
			float64(stats.Clients.Owe5),
			labels...); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeClientsTotal,
			prometheus.GaugeValue,
			float64(stats.Clients.Total),
			labels...); err == nil {
			metrics <- m
		}

		// Memory
		if m, err := prometheus.NewConstMetric(
			promDescNodeMemBuffers,
			prometheus.GaugeValue,
			float64(stats.Memory.Buffers),
			labels...); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeMemCached,
			prometheus.GaugeValue,
			float64(stats.Memory.Cached),
			labels...); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeMemFree,
			prometheus.GaugeValue,
			float64(stats.Memory.Free),
			labels...); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeMemTotal,
			prometheus.GaugeValue,
			float64(stats.Memory.Total),
			labels...); err == nil {
			metrics <- m
		}
		if m, err := prometheus.NewConstMetric(
			promDescNodeMemAvailable,
			prometheus.GaugeValue,
			float64(stats.Memory.Available),
			labels...); err == nil {
			metrics <- m
		}

		if procstat := stats.ProcStats; procstat != nil {
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUUser,
				prometheus.GaugeValue,
				float64(procstat.CPU.User),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUNice,
				prometheus.GaugeValue,
				float64(procstat.CPU.Nice),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUSystem,
				prometheus.GaugeValue,
				float64(procstat.CPU.System),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUIdle,
				prometheus.GaugeValue,
				float64(procstat.CPU.Idle),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUIOWait,
				prometheus.GaugeValue,
				float64(procstat.CPU.IOWait),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUIRQ,
				prometheus.GaugeValue,
				float64(procstat.CPU.IRQ),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatCPUSoftIRQ,
				prometheus.GaugeValue,
				float64(procstat.CPU.SoftIRQ),
				labels...); err == nil {
					metrics <- m
			}

			if m, err := prometheus.NewConstMetric(
				promDescNodeStatIntr,
				prometheus.GaugeValue,
				float64(procstat.Intr),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatContextSwitches,
				prometheus.GaugeValue,
				float64(procstat.ContextSwitches),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatSoftIRQ,
				prometheus.GaugeValue,
				float64(procstat.SoftIRQ),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeStatProcesses,
				prometheus.GaugeValue,
				float64(procstat.Processes),
				labels...); err == nil {
					metrics <- m
			}
		}
		if t := stats.Traffic.Rx; t != nil {
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficRxBytes,
				prometheus.GaugeValue,
				float64(t.Bytes),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficRxPackets,
				prometheus.GaugeValue,
				float64(t.Packets),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficRxDropped,
				prometheus.GaugeValue,
				float64(t.Dropped),
				labels...); err == nil {
					metrics <- m
			}
		}
		if t := stats.Traffic.Tx; t != nil {
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficTxBytes,
				prometheus.GaugeValue,
				float64(t.Bytes),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficTxPackets,
				prometheus.GaugeValue,
				float64(t.Packets),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficTxDropped,
				prometheus.GaugeValue,
				float64(t.Dropped),
				labels...); err == nil {
					metrics <- m
			}
		}
		if t := stats.Traffic.Forward; t != nil {
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficForwardBytes,
				prometheus.GaugeValue,
				float64(t.Bytes),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficForwardPackets,
				prometheus.GaugeValue,
				float64(t.Packets),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficForwardDropped,
				prometheus.GaugeValue,
				float64(t.Dropped),
				labels...); err == nil {
					metrics <- m
			}
		}
		if t := stats.Traffic.MgmtTx; t != nil {
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficMgmtTxBytes,
				prometheus.GaugeValue,
				float64(t.Bytes),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficMgmtTxPackets,
				prometheus.GaugeValue,
				float64(t.Packets),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficMgmtTxDropped,
				prometheus.GaugeValue,
				float64(t.Dropped),
				labels...); err == nil {
					metrics <- m
			}
		}
		if t := stats.Traffic.MgmtRx; t != nil {
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficMgmtRxBytes,
				prometheus.GaugeValue,
				float64(t.Bytes),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficMgmtRxPackets,
				prometheus.GaugeValue,
				float64(t.Packets),
				labels...); err == nil {
					metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeTrafficMgmtRxDropped,
				prometheus.GaugeValue,
				float64(t.Dropped),
				labels...); err == nil {
					metrics <- m
			}
		}


		// add a label for frequency_name
		labelIndex := len(labels)
		labelsAirtime := append(labels, "", "")
		for _, airtime := range stats.Wireless {
			labelsAirtime[labelIndex] = strconv.Itoa(int(airtime.Frequency))
			labelsAirtime[labelIndex+1] = airtime.FrequencyName()

			if m, err := prometheus.NewConstMetric(
				promDescNodeFrequency,
				prometheus.GaugeValue,
				float64(airtime.Frequency),
				labelsAirtime...); err == nil {
				metrics <- m
			}

			if m, err := prometheus.NewConstMetric(
				promDescNodeAirtimeChanUtil,
				prometheus.GaugeValue,
				float64(airtime.ChanUtil),
				labelsAirtime...); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeAirtimeTxUtil,
				prometheus.GaugeValue,
				float64(airtime.TxUtil),
				labelsAirtime...); err == nil {
				metrics <- m
			}

			if m, err := prometheus.NewConstMetric(
				promDescNodeAirtimeRxUtil,
				prometheus.GaugeValue,
				float64(airtime.RxUtil),
				labelsAirtime...); err == nil {
				metrics <- m
			}
			if m, err := prometheus.NewConstMetric(
				promDescNodeAirtimeNoise,
				prometheus.GaugeValue,
				float64(airtime.Noise),
				labelsAirtime...); err == nil {
				metrics <- m
			}
			if wireless := nodeinfo.Wireless; wireless != nil {
				if airtime.Frequency < 5000 {
					if m, err := prometheus.NewConstMetric(
						promDescNodeWirelessTxPower,
						prometheus.GaugeValue,
						float64(wireless.TxPower24),
						labelsAirtime...); err == nil {
						metrics <- m
					}
				} else {
					if m, err := prometheus.NewConstMetric(
						promDescNodeWirelessTxPower,
						prometheus.GaugeValue,
						float64(wireless.TxPower5),
						labelsAirtime...); err == nil {
						metrics <- m
					}
				}
			}
		}
	}
}
