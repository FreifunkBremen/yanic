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
	"github.com/FreifunkBremen/respond-collector/meshviewer"
)

// Node struct
type Node struct {
	Firstseen  jsontime.Time     `json:"firstseen"`
	Lastseen   jsontime.Time     `json:"lastseen"`
	Flags      *meshviewer.Flags `json:"flags,omitempty"`
	Statistics *data.Statistics  `json:"statistics"`
	Nodeinfo   *data.NodeInfo    `json:"nodeinfo"`
	Neighbours *data.Neighbours  `json:"-"`
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
			Flags: &meshviewer.Flags{
				Online:  true,
				Gateway: false,
			},
		}
		nodes.List[nodeID] = node
	}
	nodes.Unlock()

	node.Lastseen = now

	if node.Flags != nil {
		node.Flags.Online = true
	}

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
		node.Statistics = val
	}
}
func (nodes *Nodes) GetNodesMini() *meshviewer.Nodes {
	meshviewerNodes := &meshviewer.Nodes{
		Version:   1,
		List:      make(map[string]*meshviewer.Node),
		Timestamp: nodes.Timestamp,
	}
	for nodeID, _ := range nodes.List {
		meshviewerNodes.Lock()
		node, _ := meshviewerNodes.List[nodeID]

		if node == nil {
			node = &meshviewer.Node{
				Firstseen: nodes.List[nodeID].Firstseen,
				Lastseen:  nodes.List[nodeID].Lastseen,
				Flags:     nodes.List[nodeID].Flags,
				Nodeinfo:  nodes.List[nodeID].Nodeinfo,
			}
			meshviewerNodes.List[nodeID] = node
		}
		meshviewerNodes.Unlock()
		node.Statistics = &meshviewer.Statistics{
			NodeId:      nodes.List[nodeID].Statistics.NodeId,
			Clients:     nodes.List[nodeID].Statistics.Clients.Total,
			Gateway:     nodes.List[nodeID].Statistics.Gateway,
			RootFsUsage: nodes.List[nodeID].Statistics.RootFsUsage,
			LoadAverage: nodes.List[nodeID].Statistics.LoadAverage,
			Memory:      nodes.List[nodeID].Statistics.Memory,
			Uptime:      nodes.List[nodeID].Statistics.Uptime,
			Idletime:    nodes.List[nodeID].Statistics.Idletime,
			Processes:   nodes.List[nodeID].Statistics.Processes,
			MeshVpn:     nodes.List[nodeID].Statistics.MeshVpn,
			Traffic:     nodes.List[nodeID].Statistics.Traffic,
		}
	}
	return meshviewerNodes
}

// Periodically saves the cached DB to json file
func (nodes *Nodes) worker() {
	c := time.Tick(time.Second * time.Duration(nodes.config.Nodes.SaveInterval))

	for range c {
		log.Println("saving", len(nodes.List), "nodes")
		nodes.Timestamp = jsontime.Now()
		nodes.Lock()
		//
		// set node as offline (without statistics)
		for _, node := range nodes.List {
			if node.Statistics != nil && node.Lastseen.Unix()+int64(5*nodes.config.Respondd.CollectInterval) < nodes.Timestamp.Unix() {
				if node.Flags != nil {
					node.Flags.Online = false
				}
			}
		}
		// serialize nodes
		save(nodes, nodes.config.Nodes.NodesPath)
		save(nodes.GetNodesMini(), nodes.config.Nodes.NodesMiniPath)

		if path := nodes.config.Nodes.GraphsPath; path != "" {
			save(nodes.BuildGraph(), path)
		}

		nodes.Unlock()
	}
}

func (nodes *Nodes) load() {
	path := nodes.config.Nodes.NodesPath
	log.Println("loading", path)

	if filedata, err := ioutil.ReadFile(path); err == nil {
		if err := json.Unmarshal(filedata, nodes); err == nil {
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
