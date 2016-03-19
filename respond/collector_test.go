package respond

import (
	"io/ioutil"
	"net"
	"reflect"
	"testing"

	"github.com/FreifunkBremen/respond-collector/data"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)
	var decompressed *data.NodeInfo

	// callback function
	onReceive := func(addr net.UDPAddr, msg interface{}) {
		switch msg := msg.(type) {
		case *data.NodeInfo:
			decompressed = msg
		default:
			t.Error("unexpected message:", msg)
		}
	}

	collector := &Collector{
		msgType:   reflect.TypeOf(data.NodeInfo{}),
		onReceive: onReceive,
	}

	// read testdata
	compressed, err := ioutil.ReadFile("testdata/nodeinfo.flated")
	assert.Nil(err)

	collector.parse(&Response{
		Raw: compressed,
	})

	assert.NotNil(decompressed)
	assert.Equal("f81a67a5e9c1", decompressed.NodeId)
}
