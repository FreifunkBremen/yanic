package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/FreifunkBremen/respond-collector/data"
	"github.com/FreifunkBremen/respond-collector/jsontime"
)

// Node struct
type Node struct {
	Firstseen  jsontime.Time    `json:"firstseen"`
	Lastseen   jsontime.Time    `json:"lastseen"`
	Statistics *data.Statistics `json:"statistics"`
	Nodeinfo   *data.NodeInfo   `json:"nodeinfo"`
	Neighbours *data.Neighbours `json:"-"`
}

type NodeElement struct {
	NodeId string
}

// Nodes struct: cache DB of Node's structs
type Nodes struct {
	Version   int              `json:"version"`
	Timestamp jsontime.Time    `json:"timestamp"`
	List      map[string]*Node `json:"nodes"` // the current nodemap, indexed by node ID
	sync.Mutex
}

// NewNodes create Nodes structs (cache DB)
func NewNodes() *Nodes {
	nodes := &Nodes{
		Version: 2,
		List:    make(map[string]*Node),
	}

	return nodes
}

// Update a Node
func (nodes *Nodes) Update(nodeID string, res *data.ResponseData) {
	now := jsontime.Now()

	nodes.Lock()
	node, _ := nodes.List[nodeID]

	if node == nil {
		node = &Node{
			Firstseen: now,
		}
		nodes.List[nodeID] = node
	}
	nodes.Unlock()

	node.Lastseen = now

	// Update neighbours
	if val := res.Neighbours; val != nil {
		node.Neighbours = val
	}

	// Update nodeinfo
	if val := res.NodeInfo; val != nil {
		node.Nodeinfo = val
	}

	// Update statistics
	if val := res.Statistics; val != nil {
		node.Statistics = val
	}
}

// Saves the cached DB to json file periodically
func (nodes *Nodes) Saver(config *Config) {
	c := time.Tick(time.Second * time.Duration(config.Nodes.SaveInterval))

	for range c {
		log.Println("saving", len(nodes.List), "nodes")
		nodes.Timestamp = time.Now()
		nodes.Lock()
		if path := config.Nodes.NodesPath; path != "" {
			save(nodes, path)
		}
		if path := config.Nodes.GraphsPath; path != "" {
			save(nodes.BuildGraph(config.Nodes.VpnAddresses), path)
		}
		nodes.Unlock()
	}
}

func save(input interface{}, outputFile string) {
	data, err := json.Marshal(input)

	if err != nil {
		log.Panic(err)
	}

	tmpFile := outputFile + ".tmp"

	err = ioutil.WriteFile(tmpFile, data, 0644)
	if err != nil {
		log.Panic(err)
	}
	err = os.Rename(tmpFile, outputFile)
	if err != nil {
		log.Panic(err)
	}
}
