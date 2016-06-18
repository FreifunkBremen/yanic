package data

type Frequence struct {
	TxPower uint32 `json:"txpower"`
	Channel uint32 `json:"channel"`
}

type Settings struct {
	Freq24 *Frequence `json:"freq24,omitempty"`
	Freq5  *Frequence `json:"freq5,omitempty"`
}
type SwitchPort struct {
	Speed uint32 `json:"speed"`
}
