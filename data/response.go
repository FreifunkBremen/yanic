package data

// ResponseData struct
type ResponseData struct {
	Neighbours *Neighbours `json:"neighbours"`
	NodeInfo   *NodeInfo   `json:"nodeinfo"`
	Statistics *Statistics `json:"statistics"`
}
