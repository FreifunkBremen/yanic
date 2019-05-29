package prometheus

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FreifunkBremen/yanic/data"
	"github.com/FreifunkBremen/yanic/runtime"
)

func TestOutput(t *testing.T) {
	assert := assert.New(t)

	out, err := Register(map[string]interface{}{})
	assert.Error(err)
	assert.Nil(out)

	out, err = Register(map[string]interface{}{
		"path": "/tmp/prometheus.txt",
	})
	os.Remove("/tmp/prometheus.txt")
	assert.NoError(err)
	assert.NotNil(out)

	out.Save(&runtime.Nodes{
		List: map[string]*runtime.Node{
			"wasd": {
				Online: true,
				Nodeinfo: &data.Nodeinfo{
					NodeID: "wasd",
				},
				Statistics: &data.Statistics{},
			},
		},
	})
	_, err = os.Stat("/tmp/prometheus.txt")
	assert.NoError(err)
}
