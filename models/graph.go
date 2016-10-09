package models

import (
	"fmt"
	"strings"
)

type Graph struct {
	Version int `json:"version"`
	Batadv  struct {
		Directed bool         `json:"directed"`
		Graph    []string     `json:"graph"`
		Nodes    []*GraphNode `json:"nodes"`
		Links    []*GraphLink `json:"links"`
	} `json:"batadv"`
}

type GraphNode struct {
	ID     string `json:"id"`
	NodeID string `json:"node_id"`
}
type GraphLink struct {
	Source   int     `json:"source"`
	Target   int     `json:"target"`
	VPN      bool    `json:"vpn"`
	TQ       float32 `json:"tq"`
	Bidirect bool    `json:"bidirect"`
}

type GraphBuilder struct {
	macToID map[string]string      // mapping from MAC address to node id
	links   map[string]*GraphLink  // mapping from $idA-$idB to existing link
	vpn     map[string]interface{} // IDs/addresses of VPN servers
}

func (nodes *Nodes) BuildGraph() *Graph {
	builder := &GraphBuilder{
		macToID: make(map[string]string),
		links:   make(map[string]*GraphLink),
		vpn:     make(map[string]interface{}),
	}

	builder.readNodes(nodes.List)

	graph := &Graph{Version: 1}
	graph.Batadv.Directed = false
	graph.Batadv.Nodes, graph.Batadv.Links = builder.Extract()
	return graph
}

func (builder *GraphBuilder) readNodes(nodes map[string]*Node) {
	// Fill mac->id map
	for sourceID, node := range nodes {
		if nodeinfo := node.Nodeinfo; nodeinfo != nil {
			// is VPN address?
			if nodeinfo.VPN {
				builder.vpn[sourceID] = nil
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
			for sourceAddr, _ := range neighbours.LLDP {
				builder.macToID[sourceAddr] = sourceID
			}
		}
	}

	// Add links
	for sourceID, node := range nodes {
		if node.Flags.Online {
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
					for targetAddress, _ := range neighbours {
						if targetID, found := builder.macToID[targetAddress]; found {
							builder.addLink(targetID, sourceID, 255)
						}
					}
				}
			}
		}
	}
}

func (builder *GraphBuilder) Extract() ([]*GraphNode, []*GraphLink) {
	links := make([]*GraphLink, len(builder.links))
	nodes := make([]*GraphNode, len(builder.macToID))
	idToIndex := make(map[string]int)

	// collect nodes and create mapping to index
	i := 0
	for mac, nodeID := range builder.macToID {
		nodes[i] = &GraphNode{
			ID:     mac,
			NodeID: nodeID,
		}
		idToIndex[nodeID] = i
		i++
	}

	// collect links
	i = 0
	for key, link := range builder.links {
		pos := strings.IndexByte(key, '-')

		link.Source = idToIndex[key[:pos]]
		link.Target = idToIndex[key[pos+1:]]
		links[i] = link
		i++
	}

	return nodes, links
}

func (builder *GraphBuilder) isVPN(ids ...string) bool {
	for _, id := range ids {
		if _, found := builder.vpn[id]; found {
			return true
		}
	}
	return false
}

func (builder *GraphBuilder) addLink(targetID string, sourceID string, linkTq int) {
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
