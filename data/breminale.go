package data

type Wireless struct {
	TxPower24 uint32 `json:"txpower24"`
	Channel24 uint32 `json:"channel24"`
	TxPower5  uint32 `json:"txpower5"`
	Channel5  uint32 `json:"channel5"`
}
type SwitchPort struct {
	Speed uint32 `json:"speed"`
}
