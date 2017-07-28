package respond

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	// read testdata
	compressed, err := ioutil.ReadFile("testdata/nodeinfo.flated")
	assert.Nil(err)

	res := &Response{
		Raw: compressed,
	}

	data, err := res.Parse()

	assert.NoError(err)
	assert.NotNil(data)

	assert.Equal("f81a67a5e9c1", data.NodeInfo.NodeID)
}
