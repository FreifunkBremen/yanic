package models

import (
	"fmt"
	"strings"
)

type Graph struct {
	Version int `json:"version"`
	Batadv  struct {
		Directed bool `json:"directed"`
		Graph []string `json:"graph"`
		Nodes []*GraphNode `json:"nodes"`
		Links []*GraphLink `json:"links"`
	} `json:"batadv"`
}

type GraphNode struct {
	ID       string  `json:"id"`
	NodeID   string  `json:"node_id"`
}
type GraphLink struct {
	Source   interface{}    `json:"source"`
	Target   interface{}    `json:"target"`
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
	for sourceId, node := range nodes {
		if nodeinfo := node.Nodeinfo; nodeinfo != nil {
			// is VPN address?
			if nodeinfo.VPN {
				builder.vpn[sourceId] = nil
			}
			for _,batinterface := range nodeinfo.Network.Mesh {
				interfaces := batinterface.Interfaces
				addresses := append(append(interfaces.Other, interfaces.Tunnel...), interfaces.Wireless...)

				for _, sourceAddress := range addresses {
					builder.macToID[sourceAddress] = sourceId
				}
			}
		}
	}

	// Add links
	for sourceId, node := range nodes {
		if neighbours := node.Neighbours; neighbours != nil {
			for _, batadvNeighbours := range neighbours.Batadv {
				for targetAddress, link := range batadvNeighbours.Neighbours {
					if targetId, found := builder.macToID[targetAddress]; found {
						builder.addLink(targetId, sourceId, link.Tq)
					}
				}
			}
		}
	}
}

func (builder *GraphBuilder) Extract() ([]*GraphNode,[]*GraphLink) {
	iNodes := 0
	iLinks := 0
	links := make([]*GraphLink, len(builder.links))
	nodes := make([]*GraphNode, len(builder.macToID))

	for mac, nodeID := range builder.macToID {
		nodes[iNodes] = &GraphNode{
			ID: mac,
			NodeID: nodeID,
		}
		iNodes += 1
	}
	for key, link := range builder.links {
		linkPart :=strings.Split(key,"-")
		both := 0
		for i,node := range nodes{
			if(linkPart[0] == node.NodeID){
				link.Source = i
				both += 1
				continue
			}
			if(linkPart[1]==node.NodeID){
				link.Target = i
				both += 1
				break
			}
		}
		if both == 2 {
			links[iLinks] = link
			iLinks += 1
		}
	}
	return  nodes, links[:iLinks]
}

func (builder *GraphBuilder) isVPN(ids ...string) bool {
	for _, id := range ids {
		if _, found := builder.vpn[id]; found {
			return true
		}
	}
	return false
}

func (builder *GraphBuilder) addLink(targetId string, sourceId string, linkTq int) {
	// Sort IDs to generate the key
	var key string
	if strings.Compare(sourceId, targetId) > 0 {
		key = fmt.Sprintf("%s-%s", sourceId, targetId)
	} else {
		key = fmt.Sprintf("%s-%s", targetId, sourceId)
	}

	var tq float32
	if linkTq > 0 {
		tq = float32(1.0 / (float32(linkTq) / 255.0))
	}

	if link, ok := builder.links[key]; !ok {
		builder.links[key] = &GraphLink{
			VPN:    builder.isVPN(sourceId, targetId),
			TQ:     tq,
		}
	} else {
		// Use lowest of both link qualities
		if tq < link.TQ {
			link.TQ = tq
		}
		link.Bidirect = true
	}
}
