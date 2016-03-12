package data

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatistics(t *testing.T) {
	assert := assert.New(t)
	obj := &Statistics{}
	testfile("statistics.json", obj)

	assert.Equal("f81a67a601ea", obj.NodeId)
	assert.Equal("52:54:00:a9:f7:6e", obj.Gateway)
	assert.Equal(float64(57861871176), obj.Traffic.Rx.Bytes)
	assert.Equal(uint32(35), obj.Clients.Total)
	assert.Equal(uint32(35), obj.Clients.Wifi)
	assert.Equal(uint32(30), obj.Clients.Wifi24)
	assert.Equal(uint32(8), obj.Clients.Wifi5)
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
