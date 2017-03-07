package runtime

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
)

func TestExpire(t *testing.T) {
	assert := assert.New(t)
	config := &Config{}
	config.Nodes.OfflineAfter.Duration = time.Minute * 10
	config.Nodes.PruneAfter.Duration = time.Hour * 24 * 6
	nodes := &Nodes{
		config: config,
		List:   make(map[string]*Node),
	}

	nodes.Update("expire", &data.ResponseData{})  // should expire
	nodes.Update("offline", &data.ResponseData{}) // should become offline
	nodes.Update("online", &data.ResponseData{})  // should stay online

	expire := nodes.List["expire"]
	expire.Lastseen = expire.Lastseen.Add((-6 * time.Hour * 24) - time.Minute)
	offline := nodes.List["offline"]
	offline.Lastseen = offline.Lastseen.Add((-6 * time.Hour * 24) + time.Minute)

	nodes.expire()

	// one expired?
	assert.Len(nodes.List, 2)
	assert.Nil(nodes.List["expire"])

	// one offline?
	assert.NotNil(nodes.List["offline"])
	assert.False(nodes.List["offline"].Online)

	// one online?
	assert.NotNil(nodes.List["online"])
	assert.True(nodes.List["online"].Online)
}

func TestLoadAndSave(t *testing.T) {
	assert := assert.New(t)

	config := &Config{}
	config.Nodes.StatePath = "testdata/nodes.json"

	nodes := NewNodes(config)
	nodes.load()

	tmpfile, _ := ioutil.TempFile("/tmp", "nodes")
	SaveJSON(nodes, tmpfile.Name())
	os.Remove(tmpfile.Name())

	assert.Len(nodes.List, 1)
}

func TestUpdateNodes(t *testing.T) {
	assert := assert.New(t)
	nodes := &Nodes{List: make(map[string]*Node)}
	assert.Len(nodes.List, 0)

	res := &data.ResponseData{
		Neighbours: &data.Neighbours{},
		Statistics: &data.Statistics{},
		NodeInfo:   &data.NodeInfo{},
	}
	nodes.Update("abcdef012345", res)

	assert.Len(nodes.List, 1)
}
