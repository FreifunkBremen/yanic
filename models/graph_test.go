package models

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/FreifunkBremen/respond-collector/data"
	"github.com/stretchr/testify/assert"
)

type TestNode struct {
	Nodeinfo   *data.NodeInfo   `json:"nodeinfo"`
	Neighbours *data.Neighbours `json:"neighbours"`
}

func TestGenerateGraph(t *testing.T) {
	assert := assert.New(t)
	nodes := testGetNodesByFile("node1.json", "node2.json", "node3.json")

	graph := nodes.BuildGraph()
	assert.NotNil(graph)
	assert.Equal(1, graph.Version, "Wrong Version")
	assert.NotNil(graph.Batadv, "no Batadv")
	assert.Equal(false, graph.Batadv.Directed, "directed batadv")
	assert.Equal(3, len(graph.Batadv.Nodes), "wrong Nodes count")
	assert.Equal(2, len(graph.Batadv.Links), "wrong Links count")
	// TODO more tests required
}

func testGetNodesByFile(files ...string) *Nodes {

	nodes := &Nodes{
		List: make(map[string]*Node),
	}

	for _, file := range files {
		nodes.List[file] = testGetNodeByFile(file)
	}

	return nodes
}

func testGetNodeByFile(filename string) *Node {
	testnode := &TestNode{}
	testfile(filename, testnode)
	return &Node{
		Nodeinfo:   testnode.Nodeinfo,
		Neighbours: testnode.Neighbours,
	}
}

func testfile(name string, obj interface{}) {
	file, err := ioutil.ReadFile("testdata/" + name)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(file, obj); err != nil {
		panic(err)
	}
}
