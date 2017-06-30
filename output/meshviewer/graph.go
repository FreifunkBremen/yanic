package meshviewer

import (
	"fmt"
	"strings"

	"github.com/FreifunkBremen/yanic/runtime"
)

// Graph a struct for all links between the nodes
type Graph struct {
	Version int `json:"version"`
	Batadv  struct {
		Directed bool         `json:"directed"`
		Graph    []string     `json:"graph"`
		Nodes    []*GraphNode `json:"nodes"`
		Links    []*GraphLink `json:"links"`
	} `json:"batadv"`
}

// GraphNode small struct of a node for the graph struct
type GraphNode struct {
	ID     string `json:"id"`
	NodeID string `json:"node_id"`
}

// GraphLink a struct  for the link between two nodes
type GraphLink struct {
	Source   int     `json:"source"`
	Target   int     `json:"target"`
	VPN      bool    `json:"vpn"`
	TQ       float32 `json:"tq"`
	Bidirect bool    `json:"bidirect"`
}

// GraphBuilder a temporaty struct during fill the graph from the node neighbours
type graphBuilder struct {
	macToID map[string]string      // mapping from MAC address to node id
	idToMac map[string]string      // mapping from node id to one MAC address
	links   map[string]*GraphLink  // mapping from $idA-$idB to existing link
	vpn     map[string]interface{} // IDs/addresses of VPN servers
}

// BuildGraph transform from nodes (Neighbours) to Graph
func BuildGraph(nodes *runtime.Nodes) *Graph {
	builder := &graphBuilder{
		macToID: make(map[string]string),
		idToMac: make(map[string]string),
		links:   make(map[string]*GraphLink),
		vpn:     make(map[string]interface{}),
	}

	builder.readNodes(nodes.List)

	graph := &Graph{Version: 1}
	graph.Batadv.Directed = false
	graph.Batadv.Nodes, graph.Batadv.Links = builder.extract()
	return graph
}

func (builder *graphBuilder) readNodes(nodes map[string]*runtime.Node) {
	// Fill mac->id map
	for sourceID, node := range nodes {
		if nodeinfo := node.Nodeinfo; nodeinfo != nil {
			// is VPN address?
			if nodeinfo.VPN {
				builder.vpn[sourceID] = nil
			}

			if len(nodeinfo.Network.Mac) > 0 {
				builder.idToMac[sourceID] = nodeinfo.Network.Mac
			}

			// Batman neighbours
			for _, batinterface := range nodeinfo.Network.Mesh {
				interfaces := batinterface.Interfaces
				addresses := append(append(interfaces.Other, interfaces.Tunnel...), interfaces.Wireless...)

				for _, sourceAddress := range addresses {
					builder.macToID[sourceAddress] = sourceID
				}
			}
		}

		// Iterate over local MAC addresses from LLDP
		if neighbours := node.Neighbours; neighbours != nil {
			for sourceAddr := range neighbours.LLDP {
				builder.macToID[sourceAddr] = sourceID
			}
		}
	}

	// Add links
	for sourceID, node := range nodes {
		if node.Online {
			if neighbours := node.Neighbours; neighbours != nil {
				// Batman neighbours
				for _, batadvNeighbours := range neighbours.Batadv {
					for targetAddress, link := range batadvNeighbours.Neighbours {
						if targetID, found := builder.macToID[targetAddress]; found {
							builder.addLink(targetID, sourceID, link.Tq)
						}
					}
				}
				// LLDP
				for _, neighbours := range neighbours.LLDP {
					for targetAddress := range neighbours {
						if targetID, found := builder.macToID[targetAddress]; found {
							builder.addLink(targetID, sourceID, 255)
						}
					}
				}
			}
		}
	}
}

type graphNodeCache struct {
	idToMac   map[string]string
	idToIndex map[string]int
	count     int
	Nodes     []*GraphNode
}

func newGraphNodeCache(idToMac map[string]string) *graphNodeCache {
	return &graphNodeCache{
		idToMac:   idToMac,
		idToIndex: make(map[string]int),
	}
}

func (gn *graphNodeCache) getIndex(nodeID string) int {
	index, ok := gn.idToIndex[nodeID]
	if !ok {
		node := &GraphNode{
			ID:     gn.idToMac[nodeID],
			NodeID: nodeID,
		}
		gn.Nodes = append(gn.Nodes, node)
		gn.idToIndex[nodeID] = gn.count
		index = gn.count
		gn.count++
	}
	return index
}

func (builder *graphBuilder) extract() ([]*GraphNode, []*GraphLink) {
	links := make([]*GraphLink, len(builder.links))
	cache := newGraphNodeCache(builder.idToMac)

	// collect links
	i := 0
	for key, link := range builder.links {
		pos := strings.IndexByte(key, '-')
		link.Source = cache.getIndex(key[:pos])
		link.Target = cache.getIndex(key[pos+1:])
		links[i] = link
		i++
	}
	return cache.Nodes, links
}

func (builder *graphBuilder) isVPN(ids ...string) bool {
	for _, id := range ids {
		if _, found := builder.vpn[id]; found {
			return true
		}
	}
	return false
}

func (builder *graphBuilder) addLink(targetID string, sourceID string, linkTq int) {
	// Sort IDs to generate the key
	var key string
	if strings.Compare(sourceID, targetID) > 0 {
		key = fmt.Sprintf("%s-%s", sourceID, targetID)
	} else {
		key = fmt.Sprintf("%s-%s", targetID, sourceID)
	}

	var tq float32
	if linkTq > 0 {
		tq = float32(1.0 / (float32(linkTq) / 255.0))
	}

	if link, ok := builder.links[key]; !ok {
		builder.links[key] = &GraphLink{
			VPN: builder.isVPN(sourceID, targetID),
			TQ:  tq,
		}
	} else {
		// Use lowest of both link qualities
		if tq < link.TQ {
			link.TQ = tq
		}
		link.Bidirect = true
	}
}
