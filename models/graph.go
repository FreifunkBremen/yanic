package models

import (
	"fmt"
	"strings"
)

type Graph struct {
	Version int `json:"version"`
	Batadv  struct {
		Links []*GraphLink `json:"links"`
	} `json:"batadv"`
}

type GraphLink struct {
	Source   string  `json:"source"`
	Target   string  `json:"target"`
	VPN      bool    `json:"vpn"`
	TQ       float32 `json:"tq"`
	Bidirect bool    `json:"bidirect"`
}

type GraphBuilder struct {
	macToID map[string]string     // mapping from MAC address to node id
	links   map[string]*GraphLink // mapping from $idA-$idB to existing link
}

func (nodes *Nodes) BuildGraph() *Graph {
	builder := &GraphBuilder{
		macToID: make(map[string]string),
		links:   make(map[string]*GraphLink),
	}

	builder.readNodes(nodes.List)

	graph := &Graph{Version: 2}
	graph.Batadv.Links = builder.Links()
	return graph
}

func (builder *GraphBuilder) readNodes(nodes map[string]*Node) {
	// Fill mac->id map
	for sourceId, node := range nodes {
		if neighbours := node.Neighbours; neighbours != nil {
			for sourceAddress, _ := range neighbours.Batadv {
				builder.macToID[sourceAddress] = sourceId
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

func (builder *GraphBuilder) Links() []*GraphLink {
	i := 0
	links := make([]*GraphLink, len(builder.links))

	for _, link := range builder.links {
		links[i] = link
		i += 1
	}
	return links
}

func (builder *GraphBuilder) addLink(targetId string, sourceId string, linkTq int) {
	// Order IDs to get generate the key
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
			Source: sourceId,
			Target: targetId,
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
