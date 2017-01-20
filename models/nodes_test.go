package models

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/FreifunkBremen/respond-collector/data"
	"github.com/stretchr/testify/assert"
)

func TestExpire(t *testing.T) {
	assert := assert.New(t)
	config := &Config{}
	config.Nodes.MaxAge = 6
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
	assert.Equal(2, len(nodes.List))
	assert.Nil(nodes.List["expire"])

	// one offline?
	assert.NotNil(nodes.List["offline"])
	assert.False(nodes.List["offline"].Flags.Online)

	// one online?
	assert.NotNil(nodes.List["online"])
	assert.True(nodes.List["online"].Flags.Online)
}

func TestLoadAndSave(t *testing.T) {
	assert := assert.New(t)

	config := &Config{}
	config.Nodes.NodesDynamicPath = "testdata/nodes.json"

	nodes := NewNodes(config)
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

func TestToInflux(t *testing.T) {
	assert := assert.New(t)

	node := Node{
		Statistics: &data.Statistics{
			NodeId:      "foobar",
			LoadAverage: 0.5,
		},
		Nodeinfo: &data.NodeInfo{
			Owner: &data.Owner{
				Contact: "nobody",
			},
		},
		Neighbours: &data.Neighbours{},
	}

	tags, fields := node.ToInflux()

	assert.Equal("foobar", tags.GetString("nodeid"))
	assert.Equal("nobody", tags.GetString("owner"))
	assert.Equal(0.5, fields["load"])
}
