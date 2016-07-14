package data

import (
	"math"
)

type Wireless struct {
	TxPower24 uint32 `json:"txpower24,omitempty"`
	Channel24 uint32 `json:"channel24,omitempty"`
	TxPower5  uint32 `json:"txpower5,omitempty"`
	Channel5  uint32 `json:"channel5,omitempty"`
}

type WirelessStatistics struct {
	Airtime24 *WirelessAirtime `json:"airtime24,omitempty"`
	Airtime5  *WirelessAirtime `json:"airtime5,omitempty"`
}

type WirelessAirtime struct {
	ChanUtil float32 // Channel utilization
	RxUtil   float32 // Receive utilization
	TxUtil   float32 // Transmit utilization

	Active_time uint64 `json:"active"`
	Busy_time   uint64 `json:"busy"`
	Rx_time     uint64 `json:"rx"`
	Tx_time     uint64 `json:"tx"`
	Noise       uint32 `json:"noise"`
	Frequency   uint32 `json:"frequency"`
}

// Calculates the utilization values in regard to the previous values
func (cur *WirelessStatistics) SetUtilization(prev *WirelessStatistics) {
	if cur.Airtime24 != nil {
		cur.Airtime24.SetUtilization(prev.Airtime24)
	}
	if cur.Airtime5 != nil {
		cur.Airtime5.SetUtilization(prev.Airtime5)
	}
}

// Calculates the utilization values in regard to the previous values
func (cur *WirelessAirtime) SetUtilization(prev *WirelessAirtime) {
	if prev == nil || cur.Active_time <= prev.Active_time {
		return
	}

	active := float64(cur.Active_time) - float64(prev.Active_time)
	busy := float64(cur.Busy_time) - float64(prev.Busy_time)
	rx := float64(cur.Tx_time) - float64(prev.Tx_time)
	tx := float64(cur.Rx_time) - float64(prev.Rx_time)

	// Calculate utilizations
	if active > 0 {
		cur.ChanUtil = float32(math.Min(100, 100*(busy+rx+tx)/active))
		cur.RxUtil = float32(math.Min(100, 100*rx/active))
		cur.TxUtil = float32(math.Min(100, 100*tx/active))
	}
}
