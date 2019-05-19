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

	// Test unmarshalable
	config.StatePath = "testdata/nodes-broken.json"
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

	assert.Panics(func() {
		SaveJSON(nodes, "/proc/a")
		// "open /proc/a.tmp: permission denied",
	})

	tmpfile, _ = ioutil.TempFile("/tmp", "nodes")
	assert.Panics(func() {
		SaveJSON(tmpfile.Name, tmpfile.Name())
		// "json: unsupported type: func() string",
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
		Nodeinfo: &data.Nodeinfo{},
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

	nodes.AddNode(&Node{Nodeinfo: &data.Nodeinfo{}})
	assert.Len(nodes.List, 0)

	nodes.AddNode(&Node{Nodeinfo: &data.Nodeinfo{NodeID: "blub"}})
	assert.Len(nodes.List, 1)
}

func TestLinksNodes(t *testing.T) {
	assert := assert.New(t)

	nodes := &Nodes{
		List:          make(map[string]*Node),
		ifaceToNodeID: make(map[string]string),
	}
	assert.Len(nodes.List, 0)

	nodes.Update("f4f26dd7a300", &data.ResponseData{
		Nodeinfo: &data.Nodeinfo{
			NodeID: "f4f26dd7a300",
			Network: data.Network{
				Mac: "f4:f2:6d:d7:a3:00",
			},
		},
	})
	nodes.Update("f4f26dd7a30a", &data.ResponseData{
		Nodeinfo: &data.Nodeinfo{
			NodeID: "f4f26dd7a30a",
			Network: data.Network{
				Mac: "f4:f2:6d:d7:a3:0a",
			},
		},
		Neighbours: &data.Neighbours{
			NodeID: "f4f26dd7a30a",
			Babel: map[string]data.BabelNeighbours{
				"vx_mesh_lan": {
					LinkLocalAddress: "fe80::2",
					Neighbours: map[string]data.BabelLink{
						"fe80::1337": {
							Cost: 26214,
						},
					},
				},
			},
		},
	})

	nodes.Update("f4f26dd7a30b", &data.ResponseData{
		Nodeinfo: &data.Nodeinfo{
			NodeID: "f4f26dd7a30b",
			Network: data.Network{
				Mesh: map[string]*data.NetworkInterface{
					"babel": {
						Interfaces: struct {
							Wireless []string `json:"wireless,omitempty"`
							Other    []string `json:"other,omitempty"`
							Tunnel   []string `json:"tunnel,omitempty"`
						}{
							Other: []string{"fe80::1337"},
						},
					},
				},
			},
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

	// no neighbours nodeid
	node := nodes.List["f4f26dd7a300"]
	assert.NotNil(node)
	links := nodes.NodeLinks(node)
	assert.Len(links, 0)

	// babel link
	node = nodes.List["f4f26dd7a30a"]
	assert.NotNil(node)
	links = nodes.NodeLinks(node)
	assert.Len(links, 1)
	link := links[0]
	assert.Equal("f4f26dd7a30a", link.SourceID)
	assert.Equal("fe80::2", link.SourceAddress)
	assert.Equal("f4f26dd7a30b", link.TargetID)
	assert.Equal("fe80::1337", link.TargetAddress)
	assert.Equal(float32(0.6), link.TQ)

	// batman link
	node = nodes.List["f4f26dd7a30b"]
	assert.NotNil(node)
	links = nodes.NodeLinks(node)
	assert.Len(links, 1)
	link = links[0]
	assert.Equal("f4f26dd7a30b", link.SourceID)
	assert.Equal("f4:f2:6d:d7:a3:0b", link.SourceAddress)
	assert.Equal("f4f26dd7a30a", link.TargetID)
	assert.Equal("f4:f2:6d:d7:a3:0a", link.TargetAddress)
	assert.Equal(float32(0.8), link.TQ)

	nodeid := nodes.GetNodeIDbyAddress("f4:f2:6d:d7:a3:0a")
	assert.Equal("f4f26dd7a30a", nodeid)
}
