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
	var decompressed *data.ResponseData

	// callback function
	onReceive := func(addr net.UDPAddr, res *data.ResponseData) {
		decompressed = res
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
	assert.NotNil(decompressed.NodeInfo)
	assert.Equal("f81a67a5e9c1", decompressed.NodeInfo.NodeId)
}
