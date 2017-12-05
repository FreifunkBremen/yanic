package runtime

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/lib/jsontime"
)

func TestExpire(t *testing.T) {
	assert := assert.New(t)
	config := &NodesConfig{}
	config.OfflineAfter.Duration = time.Minute * 10
	// to get default (100%) path of testing
	// config.PruneAfter.Duration = time.Hour * 24 * 6
	nodes := &Nodes{
		config:        config,
		List:          make(map[string]*Node),
		ifaceToNodeID: make(map[string]string),
	}

	nodes.Update("expire", &data.ResponseData{})  // should expire
	nodes.Update("offline", &data.ResponseData{}) // should become offline
	nodes.Update("online", &data.ResponseData{})  // should stay online

	expire := nodes.List["expire"]
	expire.Lastseen = expire.Lastseen.Add((-7 * time.Hour * 24) - time.Minute)
	offline := nodes.List["offline"]
	offline.Lastseen = offline.Lastseen.Add((-7 * time.Hour * 24) + time.Minute)

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

	config := &NodesConfig{}
	// not autoload without StatePath
	NewNodes(config)

	// Test unmarshalable /dev/null - autolead with StatePath
	config.StatePath = "/dev/null"
	nodes := NewNodes(config)
	// Test unopen able
	config.StatePath = "/root/nodes.json"
	nodes.load()
	// works ;)
	config.StatePath = "testdata/nodes.json"
	nodes.load()

	tmpfile, _ := ioutil.TempFile("/tmp", "nodes")
	config.StatePath = tmpfile.Name()
	nodes.save()
	os.Remove(tmpfile.Name())

	assert.PanicsWithValue("open /dev/null.tmp: permission denied", func() {
		SaveJSON(nodes, "/dev/null")
	})

	tmpfile, _ = ioutil.TempFile("/tmp", "nodes")
	assert.PanicsWithValue("json: unsupported type: func() string", func() {
		SaveJSON(tmpfile.Name, tmpfile.Name())
	})
	os.Remove(tmpfile.Name())

	//TODO how to test easy a failing renaming

	assert.Len(nodes.List, 2)
}

func TestUpdateNodes(t *testing.T) {
	assert := assert.New(t)
	nodes := &Nodes{
		List:          make(map[string]*Node),
		ifaceToNodeID: make(map[string]string),
	}
	assert.Len(nodes.List, 0)

	res := &data.ResponseData{
		Neighbours: &data.Neighbours{},
		Statistics: &data.Statistics{
			Wireless: data.WirelessStatistics{
				&data.WirelessAirtime{},
			},
		},
		NodeInfo: &data.NodeInfo{},
	}
	nodes.Update("abcdef012345", res)

	// Update wireless statistics by running SetUtilization
	nodes.Update("abcdef012345", res)

	assert.Len(nodes.List, 1)
}

func TestSelectNodes(t *testing.T) {
	assert := assert.New(t)

	config := &NodesConfig{}
	config.StatePath = "testdata/nodes.json"

	nodes := NewNodes(config)

	selectedNodes := nodes.Select(func(n *Node) bool {
		return true
	})
	assert.Len(selectedNodes, 2)

	selectedNodes = nodes.Select(func(n *Node) bool {
		return false
	})
	assert.Len(selectedNodes, 0)

	selectedNodes = nodes.Select(func(n *Node) bool {
		return n.Nodeinfo.NodeID == "f4f26dd7a30a"
	})
	assert.Len(selectedNodes, 1)
	time := jsontime.Time{}
	time.UnmarshalJSON([]byte("2017-03-10T12:12:01"))
	assert.Equal(time, selectedNodes[0].Firstseen)
}

func TestAddNode(t *testing.T) {
	assert := assert.New(t)
	nodes := NewNodes(&NodesConfig{})

	nodes.AddNode(&Node{})
	assert.Len(nodes.List, 0)

	nodes.AddNode(&Node{Nodeinfo: &data.NodeInfo{}})
	assert.Len(nodes.List, 0)

	nodes.AddNode(&Node{Nodeinfo: &data.NodeInfo{NodeID: "blub"}})
	assert.Len(nodes.List, 1)
}

func TestLinksNodes(t *testing.T) {
	assert := assert.New(t)

	nodes := &Nodes{
		List:          make(map[string]*Node),
		ifaceToNodeID: make(map[string]string),
	}
	assert.Len(nodes.List, 0)

	nodes.Update("f4f26dd7a30a", &data.ResponseData{
		NodeInfo: &data.NodeInfo{
			NodeID: "f4f26dd7a30a",
			Network: data.Network{
				Mac: "f4:f2:6d:d7:a3:0a",
			},
		},
	})

	nodes.Update("f4f26dd7a30b", &data.ResponseData{
		NodeInfo: &data.NodeInfo{
			NodeID: "f4f26dd7a30b",
		},
		Neighbours: &data.Neighbours{
			NodeID: "f4f26dd7a30b",
			Batadv: map[string]data.BatadvNeighbours{
				"f4:f2:6d:d7:a3:0b": {
					Neighbours: map[string]data.BatmanLink{
						"f4:f2:6d:d7:a3:0a": {
							Tq: 204, Lastseen: 0.42,
						},
					},
				},
			},
		},
	})

	node := nodes.List["f4f26dd7a30a"]
	assert.NotNil(node)
	links := nodes.NodeLinks(node)
	assert.Len(links, 0)

	node = nodes.List["f4f26dd7a30b"]
	assert.NotNil(node)
	links = nodes.NodeLinks(node)
	assert.Len(links, 1)
	link := links[0]
	assert.Equal(link.SourceID, "f4f26dd7a30b")
	assert.Equal(link.SourceAddress, "f4:f2:6d:d7:a3:0b")
	assert.Equal(link.TargetID, "f4f26dd7a30a")
	assert.Equal(link.TargetAddress, "f4:f2:6d:d7:a3:0a")
	assert.Equal(link.TQ, float32(0.8))

	nodeid := nodes.GetNodeIDbyAddress("f4:f2:6d:d7:a3:0a")
	assert.Equal("f4f26dd7a30a", nodeid)
}
