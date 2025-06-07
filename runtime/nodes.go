package runtime

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/bdlm/log"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/lib/jsontime"
)

// Nodes struct: cache DB of Node's structs
type Nodes struct {
	List                map[string]*Node        `json:"nodes"` // the current nodemap, indexed by node ID
	ifaceToNodeID       map[string]string       // mapping from MAC address to NodeID
	ifaceToLinkType     map[string]LinkType     // mapping from MAC address to LinkType
	ifaceToLinkProtocol map[string]LinkProtocol // mapping from MAC address to LinkProtocol
	config              *NodesConfig
	sync.RWMutex
}

// NewNodes create Nodes structs
func NewNodes(config *NodesConfig) *Nodes {
	nodes := &Nodes{
		List:                make(map[string]*Node),
		ifaceToNodeID:       make(map[string]string),
		ifaceToLinkType:     make(map[string]LinkType),
		ifaceToLinkProtocol: make(map[string]LinkProtocol),
		config:              config,
	}

	if config.StatePath != "" {
		nodes.load()
	}

	return nodes
}

// Start all services to manage Nodes
func (nodes *Nodes) Start() {
	go nodes.worker()
}

func (nodes *Nodes) AddNode(node *Node) {
	nodeinfo := node.Nodeinfo
	if nodeinfo == nil || nodeinfo.NodeID == "" {
		return
	}
	nodes.Lock()
	defer nodes.Unlock()
	nodes.List[nodeinfo.NodeID] = node
	nodes.readIfaces(nodeinfo, node.Neighbours, false)
}

// Update a Node
func (nodes *Nodes) Update(nodeID string, res *data.ResponseData) *Node {
	now := jsontime.Now()

	nodes.Lock()
	node := nodes.List[nodeID]

	if node == nil {
		node = &Node{
			Firstseen: now,
		}
		nodes.List[nodeID] = node
	}
	if res.Nodeinfo != nil {
		nodes.readIfaces(res.Nodeinfo, res.Neighbours, true)
	}
	nodes.Unlock()

	// Update wireless statistics
	if statistics := res.Statistics; statistics != nil {
		// Update channel utilization if previous statistics are present
		if node.Statistics != nil && node.Statistics.Wireless != nil && statistics.Wireless != nil {
			statistics.Wireless.SetUtilization(node.Statistics.Wireless)
		}
	}

	// Update fields
	node.Lastseen = now
	node.Online = true
	node.Neighbours = res.Neighbours
	node.Nodeinfo = res.Nodeinfo
	node.Statistics = res.Statistics
	node.CustomFields = res.CustomFields

	return node
}

// Select selects a list of nodes to be returned
func (nodes *Nodes) Select(f func(*Node) bool) []*Node {
	nodes.RLock()
	defer nodes.RUnlock()

	result := make([]*Node, 0, len(nodes.List))
	for _, node := range nodes.List {
		if f(node) {
			result = append(result, node)
		}
	}
	return result
}

func (nodes *Nodes) GetNodeIDbyAddress(addr string) string {
	return nodes.ifaceToNodeID[addr]
}

// NodeLinks returns a list of links to known neighbours
func (nodes *Nodes) NodeLinks(node *Node) (result []Link) {
	// Store link data
	neighbours := node.Neighbours
	if neighbours == nil || neighbours.NodeID == "" || !node.Online {
		return
	}

	for sourceMAC, batadv := range neighbours.Batadv {
		for neighbourMAC, link := range batadv.Neighbours {
			if neighbourID := nodes.ifaceToNodeID[neighbourMAC]; neighbourID != "" {
				neighbour, neighbourExists := nodes.List[neighbourID]

				link := Link{
					SourceID:      neighbours.NodeID,
					SourceAddress: sourceMAC,
					TargetID:      neighbourID,
					TargetAddress: neighbourMAC,
					TQ:            float32(link.TQ) / 255.0,
				}

				if neighbourExists && neighbour.Nodeinfo != nil {
					link.TargetHostname = neighbour.Nodeinfo.Hostname
				}
				if node.Nodeinfo != nil {
					link.SourceHostname = node.Nodeinfo.Hostname
				}
				if lt, ok := nodes.ifaceToLinkType[sourceMAC]; ok && lt != OtherLinkType {
					link.Type = lt
				} else if lt, ok := nodes.ifaceToLinkType[neighbourMAC]; ok {
					link.Type = lt
				}

				result = append(result, link)
			}
		}
	}
	for _, iface := range neighbours.Babel {
		for neighbourIP, link := range iface.Neighbours {
			if neighbourID := nodes.ifaceToNodeID[neighbourIP]; neighbourID != "" {
				link := Link{
					SourceID:      neighbours.NodeID,
					SourceAddress: iface.LinkLocalAddress,
					TargetID:      neighbourID,
					TargetAddress: neighbourIP,
					TQ:            1.0 - (float32(link.Cost) / 65535.0),
				}
				if lt, ok := nodes.ifaceToLinkType[iface.LinkLocalAddress]; ok && lt != OtherLinkType {
					link.Type = lt
				} else if lt, ok := nodes.ifaceToLinkType[neighbourIP]; ok {
					link.Type = lt
				}
				result = append(result, link)
			}
		}
	}
	for sourceMAC, neighmacs := range neighbours.LLDP {
		for _, neighbourMAC := range neighmacs {
			if neighbourID := nodes.ifaceToNodeID[neighbourMAC]; neighbourID != "" {
				link := Link{
					SourceID:      neighbours.NodeID,
					SourceAddress: sourceMAC,
					TargetID:      neighbourID,
					TargetAddress: neighbourMAC,
					// TODO maybe change LLDP for link quality / 100M or 1GE
					TQ: 1.0,
				}
				if lt, ok := nodes.ifaceToLinkType[sourceMAC]; ok && lt != OtherLinkType {
					link.Type = lt
				} else if lt, ok := nodes.ifaceToLinkType[neighbourMAC]; ok {
					link.Type = lt
				}
				result = append(result, link)
			}
		}
	}
	return result
}

// Periodically saves the cached DB to json file
func (nodes *Nodes) worker() {
	c := time.Tick(nodes.config.SaveInterval.Duration)

	for range c {
		nodes.expire()
		nodes.save()
	}
}

// Expires nodes and set nodes offline
func (nodes *Nodes) expire() {
	now := jsontime.Now()

	// Nodes last seen before expireAfter will be removed
	prunePeriod := nodes.config.PruneAfter.Duration
	if prunePeriod == 0 {
		prunePeriod = time.Hour * 24 * 7 // our default
	}
	pruneAfter := now.Add(-prunePeriod)

	// Nodes last seen within OfflineAfter are changed to 'offline'
	offlineAfter := now.Add(-nodes.config.OfflineAfter.Duration)

	// Locking foo
	nodes.Lock()
	defer nodes.Unlock()

	for id, node := range nodes.List {
		if node.Lastseen.Before(pruneAfter) {
			// expire
			delete(nodes.List, id)
		} else if node.Lastseen.Before(offlineAfter) {
			// set to offline
			node.Online = false
		}
	}
}

func updateIface[K string | LinkProtocol | LinkType](class string, addr string, dataMap map[string]K, value K, warning bool) {
	if oldValue := dataMap[addr]; oldValue != value {
		var empty K
		if oldValue != empty && warning {
			log.Warnf("override %s from %s to %s on %s", class, oldValue, value, addr)
		}
		dataMap[addr] = value
	}
}

// adds the nodes interface addresses to the internal map
func (nodes *Nodes) readIfaces(nodeinfo *data.Nodeinfo, neighbours *data.Neighbours, warning bool) {
	nodeID := nodeinfo.NodeID
	network := nodeinfo.Network

	if nodeID == "" {
		log.Warn("nodeID missing in nodeinfo")
		return
	}

	addresses := []string{network.Mac}

	for _, iface := range network.Mesh {
		for _, addr := range iface.Interfaces.Wireless {
			updateIface("interface-type", addr, nodes.ifaceToLinkType, WirelessLinkType, warning)
		}
		for _, addr := range iface.Interfaces.Tunnel {
			updateIface("interface-type", addr, nodes.ifaceToLinkType, TunnelLinkType, warning)
		}
		for _, addr := range iface.Interfaces.Other {
			updateIface("interface-type", addr, nodes.ifaceToLinkType, OtherLinkType, warning)
		}
		addresses = append(addresses, iface.Addresses()...)
	}

	for _, addr := range addresses {
		if addr == "" {
			continue
		}
		updateIface("nodeID", addr, nodes.ifaceToNodeID, nodeID, warning)
	}

	if neighbours == nil || neighbours.NodeID == "" {
		return
	}

	for sourceMAC, batadv := range neighbours.Batadv {
		updateIface("mesh-protocol", sourceMAC, nodes.ifaceToLinkProtocol, BatadvLinkProtocol, warning)
		for neighbourMAC := range batadv.Neighbours {
			updateIface("mesh-protocol", neighbourMAC, nodes.ifaceToLinkProtocol, BatadvLinkProtocol, warning)
		}
	}
	for _, iface := range neighbours.Babel {
		updateIface("mesh-protocol", iface.LinkLocalAddress, nodes.ifaceToLinkProtocol, BabelLinkProtocol, warning)
		for neighbourIP := range iface.Neighbours {
			updateIface("mesh-protocol", neighbourIP, nodes.ifaceToLinkProtocol, BabelLinkProtocol, warning)
		}
	}
	for portmac, neighmacs := range neighbours.LLDP {
		updateIface("mesh-protocol", portmac, nodes.ifaceToLinkProtocol, LLDPLinkProtocol, warning)
		for _, neighmac := range neighmacs {
			updateIface("mesh-protocol", neighmac, nodes.ifaceToLinkProtocol, LLDPLinkProtocol, warning)
		}
	}

}

func (nodes *Nodes) load() {
	path := nodes.config.StatePath

	if f, err := os.Open(path); err == nil { // transform data to legacy meshviewer
		if err = json.NewDecoder(f).Decode(nodes); err == nil {
			log.Infof("loaded %d nodes", len(nodes.List))

			nodes.Lock()
			for _, node := range nodes.List {
				if node.Nodeinfo != nil {
					nodes.readIfaces(node.Nodeinfo, node.Neighbours, false)
				}
			}
			nodes.Unlock()

		} else {
			log.Errorf("failed to unmarshal nodes: %s", err)
		}
	} else {
		log.Errorf("failed to load cached nodes: %s", err)
	}
}

func (nodes *Nodes) save() {
	// Locking foo
	nodes.RLock()
	defer nodes.RUnlock()

	// serialize nodes
	SaveJSON(nodes, nodes.config.StatePath)
}

// SaveJSON to path
func SaveJSON(input interface{}, outputFile string) {
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

// Save a slice of json objects as line-encoded JSON (JSONL) to a path.
func SaveJSONL(input []interface{}, outputFile string) {
	tmpFile := outputFile + ".tmp"

	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Panic(err)
	}

	for _, element := range input {
		err = json.NewEncoder(f).Encode(element)
		if err != nil {
			log.Panic(err)
		}
	}

	f.Close()
	if err := os.Rename(tmpFile, outputFile); err != nil {
		log.Panic(err)
	}
}
