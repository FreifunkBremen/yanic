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

// Nodes struct: cache DB of Node's structs
type Nodes struct {
	Version   int              `json:"version"`
	Timestamp jsontime.Time    `json:"timestamp"`
	List      map[string]*Node `json:"nodes"` // the current nodemap, indexed by node ID
	config    *Config
	sync.RWMutex
}

type GlobalStats struct {
	Nodes         uint32
	Clients       uint32
	ClientsWifi   uint32
	ClientsWifi24 uint32
	ClientsWifi5  uint32
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
func (nodes *Nodes) Update(nodeID string, res *data.ResponseData) *Node {
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

		// Update channel utilization if previous statistics are present
		if node.Statistics != nil && node.Statistics.Wireless != nil && val.Wireless != nil {
			val.Wireless.SetUtilization(node.Statistics.Wireless)
		}

		node.Statistics = val
	}

	return node
}

// GetNodesMini get meshviewer valide JSON
func (nodes *Nodes) GetNodesMini() *meshviewer.Nodes {
	meshviewerNodes := &meshviewer.Nodes{
		Version:   1,
		List:      make(map[string]*meshviewer.Node),
		Timestamp: nodes.Timestamp,
	}
	for nodeID := range nodes.List {

		node, _ := meshviewerNodes.List[nodeID]
		nodeOrigin := nodes.List[nodeID]
		if node == nil {
			node = &meshviewer.Node{
				Firstseen: nodeOrigin.Firstseen,
				Lastseen:  nodeOrigin.Lastseen,
				Flags:     nodeOrigin.Flags,
				Nodeinfo:  nodeOrigin.Nodeinfo,
			}
			meshviewerNodes.List[nodeID] = node
		}

		// Calculate Total
		total := nodeOrigin.Statistics.Clients.Total
		if total == 0 {
			total = nodeOrigin.Statistics.Clients.Wifi24 + nodeOrigin.Statistics.Clients.Wifi5
		}

		node.Statistics = &meshviewer.Statistics{
			NodeId:      nodeOrigin.Statistics.NodeId,
			Gateway:     nodeOrigin.Statistics.Gateway,
			RootFsUsage: nodeOrigin.Statistics.RootFsUsage,
			LoadAverage: nodeOrigin.Statistics.LoadAverage,
			Memory:      nodeOrigin.Statistics.Memory,
			Uptime:      nodeOrigin.Statistics.Uptime,
			Idletime:    nodeOrigin.Statistics.Idletime,
			Processes:   nodeOrigin.Statistics.Processes,
			MeshVpn:     nodeOrigin.Statistics.MeshVpn,
			Traffic:     nodeOrigin.Statistics.Traffic,
			Clients:     total,
		}
	}
	return meshviewerNodes
}

// Periodically saves the cached DB to json file
func (nodes *Nodes) worker() {
	c := time.Tick(time.Second * time.Duration(nodes.config.Nodes.SaveInterval))

	for range c {
		nodes.expire()
		nodes.save()
	}
}

// Expires nodes and set nodes offline
func (nodes *Nodes) expire() {
	nodes.Timestamp = jsontime.Now()

	// Nodes last seen before expireTime will be removed
	maxAge := nodes.config.Nodes.MaxAge
	if maxAge <= 0 {
		maxAge = 7 // our default
	}
	expireTime := nodes.Timestamp.Add(-time.Duration(maxAge) * time.Hour * 24)

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
			if node.Flags != nil {
				node.Flags.Online = false
			}
		}
	}
}

func (nodes *Nodes) save() {
	// Locking foo
	nodes.RLock()
	defer nodes.RUnlock()

	// serialize nodes
	save(nodes, nodes.config.Nodes.NodesPath)
	save(nodes.GetNodesMini(), nodes.config.Nodes.NodesMiniPath)

	if path := nodes.config.Nodes.GraphsPath; path != "" {
		save(nodes.BuildGraph(), path)
	}
}

// Returns global statistics for InfluxDB
func (nodes *Nodes) GlobalStats() (result *GlobalStats) {
	result = &GlobalStats{}
	nodes.Lock()
	for _, node := range nodes.List {
		if node.Flags.Online {
			result.Nodes += 1
			if stats := node.Statistics; stats != nil {
				result.Clients += stats.Clients.Total
				result.ClientsWifi24 += stats.Clients.Wifi24
				result.ClientsWifi5 += stats.Clients.Wifi5
				result.ClientsWifi += stats.Clients.Wifi
			}
		}
	}
	nodes.Unlock()
	return
}

// Returns fields for InfluxDB
func (stats *GlobalStats) Fields() map[string]interface{} {
	return map[string]interface{}{
		"nodes":          stats.Nodes,
		"clients.total":  stats.Clients,
		"clients.wifi":   stats.ClientsWifi,
		"clients.wifi24": stats.ClientsWifi24,
		"clients.wifi5":  stats.ClientsWifi5,
	}
}

func (nodes *Nodes) load() {
	path := nodes.config.Nodes.NodesPath
	log.Println("loading", path)

	if filedata, err := ioutil.ReadFile(path); err == nil {
		if err = json.Unmarshal(filedata, nodes); err == nil {
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
