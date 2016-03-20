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
	config    *Config
	sync.Mutex
}

// NewNodes create Nodes structs
func NewNodes(config *Config) *Nodes {
	nodes := &Nodes{
		List:   make(map[string]*Node),
		config: config,
	}

	if config.Nodes.NodesPath != "" {
		nodes.load()
	}
	go nodes.worker()

	nodes.Version = 2
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

// Periodically saves the cached DB to json file
func (nodes *Nodes) worker() {
	c := time.Tick(time.Second * 5)

	for range c {
		log.Println("saving", len(nodes.List), "nodes")
		nodes.Timestamp = jsontime.Now()
		nodes.Lock()

		// serialize nodes
		save(nodes, nodes.config.Nodes.NodesPath)

		if path := nodes.config.Nodes.GraphsPath; path != "" {
			save(nodes, path)
		}

		nodes.Unlock()
	}
}

func (nodes *Nodes) load() {
	path := nodes.config.Nodes.NodesPath
	log.Println("loading", path)

	if data, err := ioutil.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, nodes); err == nil {
			log.Println("loaded", len(nodes.List), "nodes")
		} else {
			log.Println("failed to unmarshal nodes:", err)
		}

	} else {
		log.Println("failed loading cached nodes:", err)
	}
}

func save(input interface{}, outputFile string) {
	data, err := json.Marshal(input)
	if err != nil {
		log.Panic(err)
	}

	tmpFile := outputFile + ".tmp"

	if err := ioutil.WriteFile(tmpFile, data, 0644); err != nil {
		log.Panic(err)
	}

	if err := os.Rename(tmpFile, outputFile); err != nil {
		log.Panic(err)
	}
}
