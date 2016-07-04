package data

type Wireless struct {
	TxPower24 uint32 `json:"txpower24,omitempty"`
	Channel24 uint32 `json:"channel24,omitempty"`
	TxPower5  uint32 `json:"txpower5,omitempty"`
	Channel5  uint32 `json:"channel5,omitempty"`
}
type SwitchPort struct {
	Speed uint32 `json:"speed"`
}

type WirelessStatistics struct {
	Airtime24 *WirelessAirtime `json:"airtime24,omitempty"`
	Airtime5  *WirelessAirtime `json:"airtime5,omitempty"`
}

type WirelessAirtime struct {
	Active    uint64 `json:"active"`
	Busy      uint64 `json:"busy"`
	Rx        uint64 `json:"rx"`
	Tx        uint64 `json:"tx"`
	Noise     uint32 `json:"noise"`
	Frequency uint32 `json:"frequency"`
}
