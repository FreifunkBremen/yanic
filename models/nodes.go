package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/FreifunkBremen/RespondCollector/data"
)

// Node struct
type Node struct {
	Firstseen  time.Time        `json:"firstseen"`
	Lastseen   time.Time        `json:"lastseen"`
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
	Timestamp time.Time        `json:"timestamp"`
	List      map[string]*Node `json:"nodes"` // the current nodemap, indexed by node ID
	sync.Mutex
}

// NewNodes create Nodes structs (cache DB)
func NewNodes() *Nodes {
	nodes := &Nodes{
		Version: 1,
		List:    make(map[string]*Node),
	}

	return nodes
}

// Get a Node by nodeid
func (nodes *Nodes) Get(nodeID string) *Node {
	now := time.Now()

	nodes.Lock()
	node, _ := nodes.List[nodeID]

	if node == nil {
		node = &Node{
			Firstseen:  now,
			Nodeinfo:   &data.NodeInfo{},
			Statistics: &data.Statistics{},
			Neighbours: &data.Neighbours{},
		}
		nodes.List[nodeID] = node
	}
	nodes.Unlock()

	node.Lastseen = now

	return node
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
