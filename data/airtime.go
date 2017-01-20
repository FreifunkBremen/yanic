package data

import (
	"math"
)

// Wireless struct
type Wireless struct {
	TxPower24 uint32 `json:"txpower24,omitempty"`
	Channel24 uint32 `json:"channel24,omitempty"`
	TxPower5  uint32 `json:"txpower5,omitempty"`
	Channel5  uint32 `json:"channel5,omitempty"`
}

// WirelessStatistics struct
type WirelessStatistics []*WirelessAirtime

// WirelessAirtime struct
type WirelessAirtime struct {
	ChanUtil float32 // Channel utilization
	RxUtil   float32 // Receive utilization
	TxUtil   float32 // Transmit utilization

	ActiveTime uint64 `json:"active"`
	BusyTime   uint64 `json:"busy"`
	RxTime     uint64 `json:"rx"`
	TxTime     uint64 `json:"tx"`
	Noise      uint32 `json:"noise"`
	Frequency  uint32 `json:"frequency"`
}

// FrequencyName to 11g or 11a
func (airtime WirelessAirtime) FrequencyName() string {
	if airtime.Frequency < 5000 {
		return "11g"
	}
	return "11a"
}

// SetUtilization Calculates the utilization values in regard to the previous values
func (current WirelessStatistics) SetUtilization(previous WirelessStatistics) {
	for _, c := range current {
		for _, p := range previous {
			if c.Frequency == p.Frequency {
				c.SetUtilization(p)
			}
		}
	}
}

// SetUtilization Calculates the utilization values in regard to the previous values
func (airtime *WirelessAirtime) SetUtilization(prev *WirelessAirtime) {
	if airtime.ActiveTime <= prev.ActiveTime {
		return
	}

	active := float64(airtime.ActiveTime) - float64(prev.ActiveTime)
	busy := float64(airtime.BusyTime) - float64(prev.BusyTime)
	rx := float64(airtime.TxTime) - float64(prev.TxTime)
	tx := float64(airtime.RxTime) - float64(prev.RxTime)

	// Calculate utilizations
	if active > 0 {
		airtime.ChanUtil = float32(math.Min(100, 100*(busy+rx+tx)/active))
		airtime.RxUtil = float32(math.Min(100, 100*rx/active))
		airtime.TxUtil = float32(math.Min(100, 100*tx/active))
	}
}
