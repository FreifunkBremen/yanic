package models

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/FreifunkBremen/respond-collector/data"
	"github.com/FreifunkBremen/respond-collector/jsontime"
	"github.com/FreifunkBremen/respond-collector/meshviewer"
)

// Nodes struct: cache DB of Node's structs
type Nodes struct {
	Version   int              `json:"version"`
	Timestamp jsontime.Time    `json:"timestamp"`
	List      map[string]*Node `json:"nodes"` // the current nodemap, indexed by node ID
	config    *Config
	sync.RWMutex
}

// NewNodes create Nodes structs
func NewNodes(config *Config) *Nodes {
	nodes := &Nodes{
		List:   make(map[string]*Node),
		config: config,
	}

	if config.Nodes.NodesDynamicPath != "" {
		nodes.load()
	}
	/**
	 * Version '-1' because the nodes.json would not be defined,
	 * it would be change with the change of the respondd application on gluon
	 */
	nodes.Version = -1
	return nodes
}

//Start all services to manage Nodes
func (nodes *Nodes) Start() {
	go nodes.worker()
}

// Update a Node
func (nodes *Nodes) Update(nodeID string, res *data.ResponseData) *Node {
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
	node.Flags.Online = true

	// Update neighbours
	if val := res.Neighbours; val != nil {
		node.Neighbours = val
	}

	// Update nodeinfo
	if val := res.NodeInfo; val != nil {
		node.Nodeinfo = val
		node.Flags.Gateway = val.VPN
	}

	// Update statistics
	if val := res.Statistics; val != nil {

		// Update channel utilization if previous statistics are present
		if node.Statistics != nil && node.Statistics.Wireless != nil && val.Wireless != nil {
			val.Wireless.SetUtilization(node.Statistics.Wireless)
		}

		node.Statistics = val
	}

	return node
}

// GetNodesV1 transform data to legacy meshviewer
func (nodes *Nodes) GetNodesV1() *meshviewer.NodesV1 {
	meshviewerNodes := &meshviewer.NodesV1{
		Version:   1,
		List:      make(map[string]*meshviewer.Node),
		Timestamp: nodes.Timestamp,
	}

	for nodeID := range nodes.List {
		nodeOrigin := nodes.List[nodeID]

		if nodeOrigin.Statistics == nil {
			continue
		}

		node := &meshviewer.Node{
			Firstseen: nodeOrigin.Firstseen,
			Lastseen:  nodeOrigin.Lastseen,
			Flags:     nodeOrigin.Flags,
			Nodeinfo:  nodeOrigin.Nodeinfo,
		}
		node.Statistics = meshviewer.NewStatistics(nodeOrigin.Statistics)
		meshviewerNodes.List[nodeID] = node
	}
	return meshviewerNodes
}

// GetNodesV2 transform data to modern meshviewers
func (nodes *Nodes) GetNodesV2() *meshviewer.NodesV2 {
	meshviewerNodes := &meshviewer.NodesV2{
		Version:   2,
		Timestamp: nodes.Timestamp,
	}
	for nodeID := range nodes.List {

		nodeOrigin := nodes.List[nodeID]
		if nodeOrigin.Statistics == nil {
			continue
		}
		node := &meshviewer.Node{
			Firstseen: nodeOrigin.Firstseen,
			Lastseen:  nodeOrigin.Lastseen,
			Flags:     nodeOrigin.Flags,
			Nodeinfo:  nodeOrigin.Nodeinfo,
		}
		node.Statistics = meshviewer.NewStatistics(nodeOrigin.Statistics)
		meshviewerNodes.List = append(meshviewerNodes.List, node)
	}
	return meshviewerNodes
}

// Periodically saves the cached DB to json file
func (nodes *Nodes) worker() {
	c := time.Tick(nodes.config.Nodes.SaveInterval.Duration)

	for range c {
		nodes.expire()
		nodes.save()
	}
}

// Expires nodes and set nodes offline
func (nodes *Nodes) expire() {
	nodes.Timestamp = jsontime.Now()

	// Nodes last seen before expireTime will be removed
	pruneAfter := nodes.config.Nodes.PruneAfter.Duration
	if pruneAfter == 0 {
		pruneAfter = time.Hour * 24 * 7 // our default
	}
	expireTime := nodes.Timestamp.Add(-pruneAfter)

	// Nodes last seen before offlineTime are changed to 'offline'
	offlineTime := nodes.Timestamp.Add(-time.Minute * 10)

	// Locking foo
	nodes.Lock()
	defer nodes.Unlock()

	for id, node := range nodes.List {
		if node.Lastseen.Before(expireTime) {
			// expire
			delete(nodes.List, id)
		} else if node.Lastseen.Before(offlineTime) {
			// set to offline
			node.Flags.Online = false
		}
	}
}

func (nodes *Nodes) load() {
	path := nodes.config.Nodes.NodesDynamicPath

	if f, err := os.Open(path); err == nil { // transform data to legacy meshviewer
		if err = json.NewDecoder(f).Decode(nodes); err == nil {
			log.Println("loaded", len(nodes.List), "nodes")
		} else {
			log.Println("failed to unmarshal nodes:", err)
		}
	} else {
		log.Println("failed to load cached nodes:", err)
	}
}

func (nodes *Nodes) save() {
	// Locking foo
	nodes.RLock()
	defer nodes.RUnlock()

	// serialize nodes
	save(nodes, nodes.config.Nodes.NodesDynamicPath)

	if path := nodes.config.Nodes.NodesPath; path != "" {
		version := nodes.config.Nodes.NodesVersion
		switch version {
		case 1:
			save(nodes.GetNodesV1(), path)
		case 2:
			save(nodes.GetNodesV2(), path)
		default:
			log.Panicf("invalid nodes version: %d", version)
		}
	}

	if path := nodes.config.Nodes.GraphPath; path != "" {
		save(nodes.BuildGraph(), path)
	}
}

func save(input interface{}, outputFile string) {
	tmpFile := outputFile + ".tmp"

	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(err)
	}

	err = json.NewEncoder(f).Encode(input)
	if err != nil {
		log.Panic(err)
	}

	f.Close()
	if err := os.Rename(tmpFile, outputFile); err != nil {
		log.Panic(err)
	}
}
