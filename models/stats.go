package models

// CounterMap to manage multiple values
type CounterMap map[string]uint32

// GlobalStats struct
type GlobalStats struct {
	Clients       uint32
	ClientsWifi   uint32
	ClientsWifi24 uint32
	ClientsWifi5  uint32
	Gateways      uint32
	Nodes         uint32

	Firmwares CounterMap
	Models    CounterMap
}

//NewGlobalStats returns global statistics for InfluxDB
func NewGlobalStats(nodes *Nodes) (result *GlobalStats) {
	result = &GlobalStats{
		Firmwares: make(CounterMap),
		Models:    make(CounterMap),
	}

	nodes.Lock()
	for _, node := range nodes.List {
		if node.Flags.Online {
			result.Nodes++
			if stats := node.Statistics; stats != nil {
				result.Clients += stats.Clients.Total
				result.ClientsWifi24 += stats.Clients.Wifi24
				result.ClientsWifi5 += stats.Clients.Wifi5
				result.ClientsWifi += stats.Clients.Wifi
			}
			if node.Flags.Gateway {
				result.Gateways++
			}
			if info := node.Nodeinfo; info != nil {
				result.Models.Increment(info.Hardware.Model)
				result.Firmwares.Increment(info.Software.Firmware.Release)
			}
		}
	}
	nodes.Unlock()
	return
}

// Increment counter in the map by one
// if the value is not empty
func (m CounterMap) Increment(key string) {
	if key != "" {
		val := m[key]
		m[key] = val + 1
	}
}

// Fields returns fields for InfluxDB
func (stats *GlobalStats) Fields() map[string]interface{} {
	return map[string]interface{}{
		"nodes":          stats.Nodes,
		"gateways":       stats.Gateways,
		"clients.total":  stats.Clients,
		"clients.wifi":   stats.ClientsWifi,
		"clients.wifi24": stats.ClientsWifi24,
		"clients.wifi5":  stats.ClientsWifi5,
	}
}
