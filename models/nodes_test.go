package models

import (
	"io/ioutil"
	"os"
	"testing"

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
