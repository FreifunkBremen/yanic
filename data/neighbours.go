package data

type Neighbours struct {
	Batadv map[string]BatadvNeighbours `json:"batadv"`
	LLDP   map[string]LLDPNeighbours   `json:"lldp"`
	//WifiNeighbours map[string]WifiNeighbours   `json:"wifi"`
	NodeId string `json:"node_id"`
}

type WifiLink struct {
	Inactive int `json:"inactive"`
	Noise    int `json:"nois"`
	Signal   int `json:"signal"`
}

type BatmanLink struct {
	Lastseen float64 `json:"lastseen"`
	Tq       int     `json:"tq"`
}

type LLDPLink struct {
	Name        string `json:"name"`
	Description string `json:"descr"`
}

type BatadvNeighbours struct {
	Neighbours map[string]BatmanLink `json:"neighbours"`
}

type WifiNeighbours struct {
	Neighbours map[string]WifiLink `json:"neighbours"`
}

type LLDPNeighbours map[string]LLDPLink
