package runtime

const (
	DISABLED_AUTOUPDATER = "disabled"
	GLOBAL_SITE          = "global"
	GLOBAL_DOMAIN        = "global"
)

// CounterMap to manage multiple values
type CounterMap map[string]uint32

// GlobalStats struct
type GlobalStats struct {
	Clients       uint32
	ClientsWifi   uint32
	ClientsWifi24 uint32
	ClientsWifi5  uint32
	ClientsOwe    uint32
	ClientsOwe24  uint32
	ClientsOwe5   uint32
	Gateways      uint32
	Nodes         uint32

	Firmwares   CounterMap
	Models      CounterMap
	Autoupdater CounterMap
}

//NewGlobalStats returns global statistics for InfluxDB
func NewGlobalStats(nodes *Nodes, sitesDomains map[string][]string) (result map[string]map[string]*GlobalStats) {
	result = make(map[string]map[string]*GlobalStats)

	result[GLOBAL_SITE] = make(map[string]*GlobalStats)
	result[GLOBAL_SITE][GLOBAL_DOMAIN] = &GlobalStats{
		Firmwares:   make(CounterMap),
		Models:      make(CounterMap),
		Autoupdater: make(CounterMap),
	}

	for site, domains := range sitesDomains {
		result[site] = make(map[string]*GlobalStats)
		result[site][GLOBAL_DOMAIN] = &GlobalStats{
			Firmwares:   make(CounterMap),
			Models:      make(CounterMap),
			Autoupdater: make(CounterMap),
		}
		for _, domain := range domains {
			result[site][domain] = &GlobalStats{
				Firmwares:   make(CounterMap),
				Models:      make(CounterMap),
				Autoupdater: make(CounterMap),
			}
		}
	}

	nodes.RLock()
	for _, node := range nodes.List {
		if node.Online {
			result[GLOBAL_SITE][GLOBAL_DOMAIN].Add(node)

			if info := node.Nodeinfo; info != nil {
				site := info.System.SiteCode
				domain := info.System.DomainCode
				if _, ok := result[site]; ok {
					result[site][GLOBAL_DOMAIN].Add(node)
					if _, ok := result[site][domain]; ok {
						result[site][domain].Add(node)
					}
				}
			}
		}
	}
	nodes.RUnlock()
	return
}

// Add values to GlobalStats
// if node is online
func (s *GlobalStats) Add(node *Node) {
	s.Nodes++
	if stats := node.Statistics; stats != nil {
		s.Clients += stats.Clients.Total
		s.ClientsWifi24 += stats.Clients.Wifi24
		s.ClientsWifi5 += stats.Clients.Wifi5
		s.ClientsWifi += stats.Clients.Wifi
		s.ClientsOwe24 += stats.Clients.Owe24
		s.ClientsOwe5 += stats.Clients.Owe5
		s.ClientsOwe += stats.Clients.Owe
	}
	if node.IsGateway() {
		s.Gateways++
	}
	if info := node.Nodeinfo; info != nil {
		s.Models.Increment(info.Hardware.Model)
		s.Firmwares.Increment(info.Software.Firmware.Release)
		if info.Software.Autoupdater.Enabled {
			s.Autoupdater.Increment(info.Software.Autoupdater.Branch)
		} else {
			s.Autoupdater.Increment(DISABLED_AUTOUPDATER)
		}
	}
}

// Increment counter in the map by one
// if the value is not empty
func (m CounterMap) Increment(key string) {
	if key != "" {
		m[key]++
	}
}
