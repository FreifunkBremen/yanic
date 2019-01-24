package data

// ResponseData struct
type ResponseData struct {
	Nodeinfo   *Nodeinfo   `json:"nodeinfo" toml:"nodeinfo"`
	Statistics *Statistics `json:"statistics" toml:"statistics"`
	Neighbours *Neighbours `json:"neighbours" toml:"neighbours"`
	CustomFields map[string]interface{} `json:"-"`
}
