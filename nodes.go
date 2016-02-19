package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type Node struct {
	Firstseen  time.Time   `json:"firstseen"`
	Lastseen   time.Time   `json:"lastseen"`
	Statistics interface{} `json:"statistics"`
	Nodeinfo   interface{} `json:"nodeinfo"`
	Neighbours interface{} `json:"neighbours"`
}

type Nodes struct {
	Version   int              `json:"version"`
	Timestamp time.Time        `json:"timestamp"`
	List      map[string]*Node `json:"nodes"` // the current nodemap
	sync.Mutex
}

func NewNodes() *Nodes {
	nodes := &Nodes{
		Version: 1,
		List:    make(map[string]*Node),
	}

	go nodes.saver()

	return nodes
}

func (nodes *Nodes) get(nodeId string) *Node {
	now := time.Now()

	nodes.Lock()
	node, _ := nodes.List[nodeId]

	if node == nil {
		node = &Node{
			Firstseen: now,
		}
		nodes.List[nodeId] = node
	}
	nodes.Unlock()

	node.Lastseen = now

	return node
}

func (nodes *Nodes) saver() {
	c := time.Tick(saveInterval)

	for range c {
		nodes.save()
	}
}

func (nodes *Nodes) save() {
	nodes.Timestamp = time.Now()

	nodes.Lock()
	data, err := json.Marshal(nodes)
	nodes.Unlock()

	if err !=nil{
		log.Panic(err)
	}
	log.Println("saving", len(nodes.List), "nodes")

	tmpFile := outputFile + ".tmp"

	err = ioutil.WriteFile(tmpFile, data, 0644)
	if err !=nil{
		log.Panic(err)
	}
	err = os.Rename(tmpFile, outputFile)
	if err !=nil{
		log.Panic(err)
	}
}
