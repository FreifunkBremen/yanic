package data

// Neighbours struct
type Neighbours struct {
	Batadv map[string]BatadvNeighbours `json:"batadv"`
	Babel  map[string]BabelNeighbours  `json:"babel"`
	LLDP   map[string]LLDPNeighbours   `json:"lldp"`
	//WifiNeighbours map[string]WifiNeighbours   `json:"wifi"`
	NodeID string `json:"node_id"`
}

// WifiLink struct
type WifiLink struct {
	Inactive int `json:"inactive"`
	Noise    int `json:"nois"`
	Signal   int `json:"signal"`
}

// BatmanLink struct
type BatmanLink struct {
	Lastseen float64 `json:"lastseen"`
	Tq       int     `json:"tq"`
}

// LLDPLink struct
type LLDPLink struct {
	Name        string `json:"name"`
	Description string `json:"descr"`
}

// BabelLink struct
type BabelLink struct {
	// How need this:
	RXCost       int `json:"rxcost"`
	TXCost       int `json:"txcost"`
	Cost         int `json:"cost"`
	Reachability int `json:"reachability"`
}

// BatadvNeighbours struct
type BatadvNeighbours struct {
	Neighbours map[string]BatmanLink `json:"neighbours"`
}

// BabelNeighbours struct
type BabelNeighbours struct {
	Protocol         string               `json:"protocol"`
	LinkLocalAddress string               `json:"ll-addr"`
	Neighbours       map[string]BabelLink `json:"neighbours"`
}

// WifiNeighbours struct
type WifiNeighbours struct {
	Neighbours map[string]WifiLink `json:"neighbours"`
}

// LLDPNeighbours struct
type LLDPNeighbours map[string]LLDPLink
