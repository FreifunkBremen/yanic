package models

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/FreifunkBremen/respond-collector/data"
	"github.com/stretchr/testify/assert"
)

func TestLoadAndSave(t *testing.T) {
	assert := assert.New(t)

	config := &Config{}
	config.Nodes.NodesPath = "testdata/nodes.json"

	nodes := &Nodes{config: config}
	nodes.load()

	tmpfile, _ := ioutil.TempFile("/tmp", "nodes")
	save(nodes, tmpfile.Name())
	os.Remove(tmpfile.Name())

	assert.Equal(1, len(nodes.List))
}

func TestUpdateNodes(t *testing.T) {
	assert := assert.New(t)
	nodes := &Nodes{List: make(map[string]*Node)}
	assert.Equal(0, len(nodes.List))

	res := &data.ResponseData{
		Neighbours: &data.Neighbours{},
		Statistics: &data.Statistics{},
		NodeInfo:   &data.NodeInfo{},
	}
	nodes.Update("abcdef012345", res)

	assert.Equal(1, len(nodes.List))
}
