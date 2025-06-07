package main

import (
	"sync"

	meshviewerFFRGB "github.com/FreifunkBremen/yanic/output/meshviewer-ffrgb"
)

type Node struct {
	NodeID    string   `json:"node_id"`
	Hostname  string   `json:"hostname"`
	Addresses []string `json:"addresses"`
}

type Status struct {
	Error             string  `json:"error,omitempty"`
	NodesCount        int     `json:"nodes_count"`
	NodesOfflineCount int     `json:"nodes_offline_count"`
	NodesCrashed      []*Node `json:"nodes_crashed"`
	sync.Mutex
}

func (s *Status) AddNode(node *meshviewerFFRGB.Node) {
	s.Lock()
	s.NodesCrashed = append(s.NodesCrashed, &Node{
		NodeID:    node.NodeID,
		Hostname:  node.Hostname,
		Addresses: node.Addresses,
	})
	s.Unlock()
}
